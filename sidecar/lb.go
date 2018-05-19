package sidecar

import (
	"errors"

	"dubbo-mesh/registry"
)

const (
	ElectorRandom = iota
	ElectorRR
	ElectorWRR
	ElectorLL
	ElectorLA
	ElectorDRR
)

func elector(elector int) Banlancer {
	switch elector {
	case ElectorRandom:
		return &Random{}
	case ElectorRR:
		return &RoundRobin{}
	case ElectorWRR:
		return &WeightRoundRobin{}
	case ElectorLL:
		return &LeastLatest{}
	case ElectorLA:
		return &LeastAVG{}
	case ElectorDRR:
		return &DynamicWeightRoundRobin{}
	default:
		panic(errors.New("unknown elector"))
	}
}

type Banlancer interface {
	Init(endpoint []*Endpoint)
	Elect(endpoints []*Endpoint) *Endpoint
}

type Endpoint struct {
	*registry.Endpoint
	Meter       *Meter
	good        bool
	curWeight   int
	originWight int
}

type Rtt struct {
	Endpoint *Endpoint
	Rtt      int64
	Error    error
}

type Meter struct {
	TotalCount uint64 // 处理的总数
	ErrorCount uint64 // 错误数
	Error      uint64 // 大于0，最近n次错误，=0 最近一次没有错误
	Latest     uint64 // RTT
	Max        uint64
	Min        uint64
	Total      uint64
}

// 平均RTT
func (this *Meter) Avg() uint64 {
	if this.TotalCount == 0 {
		return this.Total
	}
	return this.Total / this.TotalCount
}

// 正确处理率
func (this *Meter) Rate() float64 {
	return float64(this.TotalCount-this.ErrorCount) / float64(this.TotalCount)
}
