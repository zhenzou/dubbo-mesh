package sidecar

import (
	"sync/atomic"

	"dubbo-mesh/log"
)

type RoundRobin struct {
	baseBalancer
	index int32
	total int32
}

func (this *RoundRobin) Init(endpoints []*Endpoint) {
	this.baseBalancer.Init(endpoints)
	this.total = int32(len(endpoints))
}

// round robin
func (this *RoundRobin) Elect() *Endpoint {
	endpoint := this.endpoints[this.index]
	for !atomic.CompareAndSwapInt32(&this.index, this.index, (this.index+1)/this.total) {
		log.Debug("rr:", this.index)
	}
	return endpoint
}

// 加权轮询
type WeightRoundRobin struct {
	baseBalancer
	weights map[*Endpoint]int
	index   int
	max     int
	cw      int // 当前权重
}

func (this *WeightRoundRobin) Init(endpoints []*Endpoint) {
	this.baseBalancer.Init(endpoints)
	this.index = -1
	this.weights = make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weight := this.weight(endpoint)
		this.weights[endpoint] = weight
	}
	gcd := this.gcd(this.weights)
	for k, weight := range this.weights {
		max := weight / gcd
		this.weights[k] = max
		if max > this.max {
			this.max = max
		}
	}
}

func (this *WeightRoundRobin) Elect() *Endpoint {
	endpoints := this.endpoints
	for {
		this.index = (this.index + 1) % len(endpoints)
		if this.index == 0 {
			this.cw = this.cw - 1
			if this.cw <= 0 {
				this.cw = this.max
			}
		}
		end := endpoints[this.index]
		if this.weights[end] >= this.cw {
			return end
		}
	}
}

const (
	SwitchCount = 20000
)

type WeightLeastLatestAvg struct {
	baseBalancer
	wrr   WeightRoundRobin
	lla   LeastLatestAvg
	count int32
}

func (this *WeightLeastLatestAvg) Init(endpoints []*Endpoint) {
	this.wrr.Init(endpoints)
	this.lla.Init(endpoints)
}

// 用wrr预热
func (this *WeightLeastLatestAvg) Elect() *Endpoint {
	count := atomic.AddInt32(&this.count, 1)
	if count < SwitchCount {
		return this.wrr.Elect()
	} else {
		return this.lla.Elect()
	}
}
