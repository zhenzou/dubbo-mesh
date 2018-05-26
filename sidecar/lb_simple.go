package sidecar

import (
	"math"

	"dubbo-mesh/util"
)

// 下面的两种会导致性能最好的负担太重
// 最近响应时间最短
// 不太行
type LeastLatest struct {
}

func (this *LeastLatest) Init(endpoint []*Endpoint) {
	// do nothing
}

func (this *LeastLatest) Elect(endpoints []*Endpoint) *Endpoint {
	var result *Endpoint
	var min uint64 = math.MaxUint64
	for _, endpoint := range endpoints {
		if endpoint.Meter.Latest < min {
			min = endpoint.Meter.Latest
			result = endpoint
		}
	}
	return result
}

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
		if act := endpoint.Active; act < min {
			min = act
			result = endpoint
		}
	}
	return result
}

type WeightLeastActive struct {
	weights map[*Endpoint]int32
}

func (this *WeightLeastActive) Init(endpoints []*Endpoint) {
	this.weights = make(map[*Endpoint]int32)
	for _, endpoint := range endpoints {
		weight := this.calculateWrr(endpoint)
		this.weights[endpoint] = weight
	}
	gcd := this.weightGcd()
	for k, weight := range this.weights {
		max := weight / gcd
		this.weights[k] = max
	}
}

func (r *WeightLeastActive) weightGcd() int32 {
	var divisor int32 = -1
	for _, s := range r.weights {
		if divisor == -1 {
			divisor = s
		} else {
			divisor = util.Gcd32(divisor, s)
		}
	}
	return divisor
}

// 简单的计算权重，暂时 就把内存做为权重
func (this *WeightLeastActive) calculateWrr(status *Endpoint) int32 {
	return int32(status.System.TotalMemory)
}

func (this *WeightLeastActive) weight(endpoint *Endpoint) int32 {

	return endpoint.Active / this.weights[endpoint]
}

func (this *WeightLeastActive) Elect(endpoints []*Endpoint) *Endpoint {
	var result *Endpoint
	var min int32 = math.MaxInt32
	for _, endpoint := range endpoints {
		if cur := this.weight(endpoint); cur < min {
			min = cur
			result = endpoint
		}
	}
	return result
}
