package sidecar

import (
	"math"

	"dubbo-mesh/util"
)

// 平均响应时间最短
// 不太行
type LeastAVG struct {
}

func (this *LeastAVG) Init(endpoint []*Endpoint) {
	// do nothing
}

func (this *LeastAVG) Elect(endpoints []*Endpoint) *Endpoint {
	var result *Endpoint
	var min uint64 = math.MaxUint64
	for _, endpoint := range endpoints {
		if avg := endpoint.Meter.Avg(); avg < min {
			min = avg
			result = endpoint
		}
	}
	return result
}

type LeastActive struct {
}

func (this *LeastActive) Init(endpoints []*Endpoint) {
	// do nothing
}

func (this *LeastActive) Elect(endpoints []*Endpoint) *Endpoint {
	var result *Endpoint
	var min int32 = math.MaxInt32
	for _, endpoint := range endpoints {
		if act := endpoint.Meter.Active; act < min {
			min = act
			result = endpoint
		}
	}
	return result
}

type WeightLeastLatestAvg struct {
	weights map[*Endpoint]int
	next    *Endpoint
}

func (this *WeightLeastLatestAvg) Init(endpoints []*Endpoint) {
	this.weights = make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weight := this.calculateWrr(endpoint)
		this.weights[endpoint] = weight
	}
	gcd := this.weightGcd()
	total := 0
	for k, w := range this.weights {
		weight := w / gcd
		total += weight
		this.weights[k] = weight
	}
	this.next = endpoints[0]
}

func (r *WeightLeastLatestAvg) weightGcd() int {
	divisor := -1
	for _, s := range r.weights {
		if divisor == -1 {
			divisor = s
		} else {
			divisor = util.Gcd(divisor, s)
		}
	}
	return divisor
}

// 简单的计算权重，暂时 就把内存做为权重
func (this *WeightLeastLatestAvg) calculateWrr(status *Endpoint) int {
	return status.System.Memory
}

func (this *WeightLeastLatestAvg) weight(endpoint *Endpoint) int {
	return int(endpoint.Meter.Avg())
}

func (this *WeightLeastLatestAvg) Elect(endpoints []*Endpoint) *Endpoint {
	min := math.MaxInt32
	result := this.next
	for _, endpoint := range endpoints {
		if *endpoint == *result {
			continue
		}
		if cur := this.weight(endpoint); cur < min {
			min = cur
			this.next = endpoint
		}
	}
	return result
}
