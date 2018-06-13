package sidecar

import (
	"math"
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

// 最小最近平均latency
type LeastLatestAvg struct {
	next *Endpoint
}

func (this *LeastLatestAvg) Init(endpoints []*Endpoint) {
	this.next = endpoints[0]
}

func (this *LeastLatestAvg) Elect(endpoints []*Endpoint) *Endpoint {
	min := math.MaxInt32
	result := this.next
	for _, endpoint := range endpoints {
		if *endpoint == *result {
			continue
		}
		if cur := int(endpoint.Meter.Avg()); cur < min {
			min = cur
			this.next = endpoint
		}
	}
	return result
}
