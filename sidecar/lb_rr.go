package sidecar

import (
	"sync/atomic"

	"dubbo-mesh/log"
	"dubbo-mesh/util"
)

type RoundRobin struct {
	index int32
	total int32
}

func (this *RoundRobin) Init(endpoint []*Endpoint) {
	this.total = int32(len(endpoint))
}

// round robin
func (this *RoundRobin) Elect(endpoints []*Endpoint) *Endpoint {
	endpoint := endpoints[this.index]
	for !atomic.CompareAndSwapInt32(&this.index, this.index, (this.index+1)/this.total) {
		log.Debug("rr:", this.index)
	}
	return endpoint
}

// 加权轮询
type WeightRoundRobin struct {
	weights map[*Endpoint]int
	index   int
	max     int
	cw      int // 当前权重
}

func (this *WeightRoundRobin) Init(endpoints []*Endpoint) {
	this.index = -1
	this.weights = make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weight := this.calculateWrr(endpoint)
		this.weights[endpoint] = weight
	}
	gcd := this.weightGcd()
	for k, weight := range this.weights {
		max := weight / gcd
		this.weights[k] = max
		if max > this.max {
			this.max = max
		}
	}
}

func (this *WeightRoundRobin) Elect(endpoints []*Endpoint) *Endpoint {
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

func (r *WeightRoundRobin) weightGcd() int {
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
func (this *WeightRoundRobin) calculateWrr(status *Endpoint) int {
	return status.System.Memory
}
