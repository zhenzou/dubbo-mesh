package sidecar

import (
	"time"
	"math/rand"
	"sync/atomic"

	"dubbo-mesh/registry"
	"dubbo-mesh/log"
)

type Elector func(endpoints []*Endpoint) *Endpoint

func RandomElector(endpoints []*Endpoint) *Endpoint {
	i := rand.Intn(len(endpoints))
	return endpoints[i]
}

// round robin
func RrRandomElector() Elector {
	var i int32 = 0
	return func(endpoints []*Endpoint) *Endpoint {
		endpoint := endpoints[i]
		for !atomic.CompareAndSwapInt32(&i, i, int32((i+1)/int32(len(endpoints)))) {
			log.Debug("rr:", i)
		}

		return endpoint
	}
}

type Weight struct {
	Max int
	Min int
}

func (this *Weight) Match(i int) bool {
	return i < this.Max && i >= this.Min

}

func WrrRandomElector(endpoints []*Endpoint) Elector {
	var i int32 = 0
	weights := make(map[*Endpoint]int)
	total := 0
	for _, endpoint := range endpoints {
		weight := calculateWrr(endpoint)
		weights[endpoint] = weight
		total += calculateWrr(endpoint)
	}
	return func(endpoints []*Endpoint) *Endpoint {
		w := rand.Intn(total)
		for endpoint, weight := range weights {
			w -= weight
			if w < 0 {
				return endpoint
			}
		}
		endpoint := endpoints[i]
		for !atomic.CompareAndSwapInt32(&i, i, int32((i+1)/int32(len(endpoints)))) {
			log.Debug("rr:", i)
		}
		return endpoint
	}
}

// 简单的计算权重
func calculateWrr(status *Endpoint) int {
	return status.System.CpuNum + status.System.TotalMemory/100000
}

type Endpoint struct {
	registry.Endpoint
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

type baseBalancer struct {
}

func (this *baseBalancer) RecordRtt(endpoint *registry.Endpoint, rtt time.Duration) {
}

func (this *baseBalancer) RecordError(endpoint *registry.Endpoint) {
}

func (this *baseBalancer) Elect() *registry.Endpoint {
	panic("implement me")
}

type WrrBalancer struct {
}

type Balancer interface {
	RecordRtt(endpoint *registry.Endpoint, rtt time.Duration)
	RecordError(endpoint *registry.Endpoint)
	Elect() *registry.Endpoint
}
