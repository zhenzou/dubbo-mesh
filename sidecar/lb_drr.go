package sidecar

import (
	"time"
	"sync"
	"math"
	"sort"

	"dubbo-mesh/util"
	"dubbo-mesh/log"
)

// 动态加权轮询

// 动态权重变化
// 基本从oxy复制来的
type DynamicWeightRoundRobin struct {
	*WeightRoundRobin
	ratings   []float64
	timer     time.Time
	mtx       *sync.Mutex
	endpoints []*Endpoint
}

const (
	splitThreshold = 1.5 // drr用到的分隔比例
	// This is the maximum weight that handler will set for the server
	FSMMaxWeight = 4096
	// Multiplier for the server weight
	FSMGrowFactor = 4
)

func (this *DynamicWeightRoundRobin) Init(endpoints []*Endpoint) {
	this.WeightRoundRobin = &WeightRoundRobin{}
	this.WeightRoundRobin.Init(endpoints)
	this.mtx = &sync.Mutex{}
	this.ratings = make([]float64, len(endpoints))
	this.endpoints = endpoints
	go this.cron()
}

func (this *DynamicWeightRoundRobin) Elect(endpoints []*Endpoint) *Endpoint {
	endpoint := this.WeightRoundRobin.Elect(endpoints)
	log.Debugf("weight:%d,avg:%d ", endpoint.curWeight, endpoint.Meter.Avg())
	return endpoint
}

func (this *DynamicWeightRoundRobin) cron() {
	tick := time.NewTicker(1 * time.Second)
	for _ = range tick.C {
		this.adjustWeight(this.endpoints)
	}
}

// 动态的计算权重，考虑系统配置和运行状态
func (this *DynamicWeightRoundRobin) adjustWeight(endpoints []*Endpoint) {
	if len(this.weights) < 2 {
		return
	}
	this.mtx.Lock()
	defer this.mtx.Unlock()

	if this.mark(endpoints) {
		if this.setMarkedWeights(endpoints) {
			this.setTimer()
		}
	} else { // No servers that are different by their quality, so converge weights
		if this.convergeWeights(endpoints) {
			this.setTimer()
		}
	}
}

func (rb *DynamicWeightRoundRobin) convergeWeights(endpoints []*Endpoint) bool {
	// If we have previously changed servers try to restore weights to the original state
	changed := false
	for _, s := range endpoints {
		if s.originWight == s.curWeight {
			continue
		}
		changed = true
		newWeight := decrease(s.originWight, s.curWeight)
		s.curWeight = newWeight
	}
	if !changed {
		return false
	}
	rb.normalizeWeights(endpoints)
	rb.applyWeights(endpoints)
	return true
}

func (rb *DynamicWeightRoundRobin) setTimer() {
	rb.timer = rb.timer.Add(10 * time.Second)
}

func (rb *DynamicWeightRoundRobin) setMarkedWeights(endpoints []*Endpoint) bool {
	changed := false
	// Increase weights on servers marked as good
	for _, srv := range endpoints {
		if srv.good {
			weight := increase(srv.curWeight)
			if weight <= FSMMaxWeight {
				srv.curWeight = weight
				changed = true
			}
		}
	}
	if changed {
		rb.normalizeWeights(endpoints)
		rb.applyWeights(endpoints)
		return true
	}
	return false
}

func (rb *DynamicWeightRoundRobin) applyWeights(endpoints []*Endpoint) {
	for _, endpoint := range endpoints {
		rb.weights[endpoint] = endpoint.curWeight
	}
}

func (rb *DynamicWeightRoundRobin) normalizeWeights(endpoints []*Endpoint) {
	gcd := rb.weightsGcd(endpoints)
	if gcd <= 1 {
		return
	}
	for _, e := range endpoints {
		e.curWeight = e.curWeight / gcd
	}
}

func (rb *DynamicWeightRoundRobin) weightsGcd(endpoints []*Endpoint) int {
	divisor := -1
	for _, w := range endpoints {
		if divisor == -1 {
			divisor = w.curWeight
		} else {
			divisor = util.Gcd(divisor, w.curWeight)
		}
	}
	return divisor
}

// 动态的计算权重，考虑系统配置和运行状态
func (this *DynamicWeightRoundRobin) mark(endpoints []*Endpoint) bool {
	for i, srv := range endpoints {
		this.ratings[i] = srv.Meter.Rate()
	}
	g, b := SplitFloat64(splitThreshold, 0, this.ratings)

	for i, srv := range endpoints {
		if g[this.ratings[i]] {
			srv.good = true
		} else {
			srv.good = false
		}
	}
	return len(g) != 0 && len(b) != 0

}

func decrease(target, current int) int {
	adjusted := current / FSMGrowFactor
	if adjusted < target {
		return target
	} else {
		return adjusted
	}
}

func increase(weight int) int {
	return weight * FSMGrowFactor
}

func SplitFloat64(threshold, sentinel float64, values []float64) (good map[float64]bool, bad map[float64]bool) {
	good, bad = make(map[float64]bool), make(map[float64]bool)
	var newValues []float64
	if len(values)%2 == 0 {
		newValues = make([]float64, len(values)+1)
		copy(newValues, values)
		// Add a sentinel endpoint so we can distinguish outliers better
		newValues[len(newValues)-1] = sentinel
	} else {
		newValues = values
	}

	m := median(newValues)
	mAbs := medianAbsoluteDeviation(newValues)
	for _, v := range values {
		if v > (m+mAbs)*threshold {
			bad[v] = true
		} else {
			good[v] = true
		}
	}
	return good, bad
}

func medianAbsoluteDeviation(values []float64) float64 {
	m := median(values)
	distances := make([]float64, len(values))
	for i, v := range values {
		distances[i] = math.Abs(v - m)
	}
	return median(distances)
}

func median(values []float64) float64 {
	vals := make([]float64, len(values))
	copy(vals, values)
	sort.Float64s(vals)
	l := len(vals)
	if l%2 != 0 {
		return vals[l/2]
	}
	return (vals[l/2-1] + vals[l/2]) / 2.0
}
