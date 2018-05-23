package sidecar

import (
	"math"
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

func (this *LeastActive) Init(endpoint []*Endpoint) {
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
