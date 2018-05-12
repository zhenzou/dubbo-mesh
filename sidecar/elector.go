package sidecar

import (
	"time"
	"math/rand"
	"sync/atomic"
	"errors"

	"dubbo-mesh/registry"
	"dubbo-mesh/log"
)

const (
	ElectorRandom = iota
	ElectorRR
	ElectorWRR
)

func elector(elector int) Elector {
	switch elector {
	case ElectorRandom:
		return &Random{}
	case ElectorRR:
		return &RoundRobin{}
	case ElectorWRR:
		return &WrrRandom{}
	default:
		panic(errors.New("unknown elector"))
	}
}

type Elector interface {
	Init(endpoint []*Endpoint)
	Elect(endpoints []*Endpoint) *Endpoint
}

type Random struct {
	total int
}

func (this *Random) Init(endpoint []*Endpoint) {
	this.total = len(endpoint)
}

func (this *Random) Elect(endpoints []*Endpoint) *Endpoint {
	i := rand.Intn(this.total)
	return endpoints[i]
}

type RoundRobin struct {
	index int32
}

func (this *RoundRobin) Init(endpoint []*Endpoint) {
}

// round robin
func (this *RoundRobin) Elect(endpoints []*Endpoint) *Endpoint {
	endpoint := endpoints[this.index]
	for !atomic.CompareAndSwapInt32(&this.index, this.index, int32((this.index+1)/int32(len(endpoints)))) {
		log.Debug("rr:", this.index)
	}

	return endpoint
}

type WrrRandom struct {
	weights map[*Endpoint]int
	total   int
}

func (this *WrrRandom) Init(endpoints []*Endpoint) {
	this.weights = make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weight := calculateWrr(endpoint)
		this.weights[endpoint] = weight
		this.total += calculateWrr(endpoint)
	}
}

func (this *WrrRandom) Elect(endpoints []*Endpoint) *Endpoint {
	w := rand.Intn(this.total)
	for endpoint, weight := range this.weights {
		w -= weight
		if w < 0 {
			return endpoint
		}
	}
	return nil
}

// 简单的计算权重
func calculateWrr(status *Endpoint) int {
	return status.System.CpuNum + status.System.TotalMemory/100000
}

func NewEndpoint(endpoint *registry.Endpoint) *Endpoint {
	return &Endpoint{
		Endpoint: endpoint,
		Status:   &Status{},
	}
}

type Endpoint struct {
	*registry.Endpoint
	Status *Status
}

type Status struct {
	Count  int           // 处理的总数
	Rate   int           // 处理速率
	Latest time.Duration // RTT
	Max    time.Duration
	Min    time.Duration
	Avg    time.Duration
}
