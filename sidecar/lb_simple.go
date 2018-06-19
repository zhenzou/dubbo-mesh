package sidecar

import (
	"math"
	"sync"
)

// 平均响应时间最短
// 不太行
type LeastAVG struct {
	baseBalancer
}

func (this *LeastAVG) Elect() *Endpoint {
	var result *Endpoint
	var min uint64 = math.MaxUint64
	for _, endpoint := range this.endpoints {
		if avg := endpoint.Meter.Avg(); avg < min {
			min = avg
			result = endpoint
		}
	}
	return result
}

type LeastActive struct {
	WeightRoundRobin
	next *Endpoint
}

func (this *LeastActive) Init(endpoints []*Endpoint) {
	this.WeightRoundRobin.Init(endpoints)
	this.next = endpoints[0]
}

func (this *LeastActive) Record(endpoint *Endpoint, latency uint64) {
	meter := endpoint.Meter
	meter.Mtx.Lock()
	meter.Count += 1
	meter.Total += latency
	if meter.Count >= AvgCount {
		val := meter.Queue.Remove()
		meter.Total -= val.(uint64)
	}
	meter.Queue.Add(latency)
	meter.Mtx.Unlock()
}

func (this *LeastActive) Elect() *Endpoint {
	result := this.next
	min := math.MaxInt32
	for _, endpoint := range this.endpoints {
		if cur := int(endpoint.Meter.Active) * int(endpoint.Meter.Avg()) / this.weights[endpoint]; cur < min {
			min = cur
			this.next = endpoint
		}
	}
	return result
}

// 最小最近平均latency
type LeastLatestAvg struct {
	WeightRoundRobin
	next *Endpoint
	mtx  sync.Mutex
}

func (this *LeastLatestAvg) Init(endpoints []*Endpoint) {
	this.WeightRoundRobin.Init(endpoints)
	this.next = endpoints[0]
}

const (
	AvgCount = 16
)

func (this *LeastLatestAvg) Record(endpoint *Endpoint, latency uint64) {
	meter := endpoint.Meter
	meter.Mtx.Lock()
	meter.Count += 1
	meter.Total += latency
	if meter.Count >= AvgCount {
		val := meter.Queue.Remove()
		meter.Total -= val.(uint64)
	}
	meter.Queue.Add(latency)
	meter.Mtx.Unlock()
}

func (this *LeastLatestAvg) Elect() *Endpoint {
	min := math.MaxInt32
	result := this.next
	endpoints := this.endpoints
	for _, endpoint := range endpoints {
		if endpoint == result {
			continue
		}
		meter := endpoint.Meter
		if cur := int(meter.Avg()); cur < min {
			min = cur
			this.next = endpoint
		}
	}
	return result
}
