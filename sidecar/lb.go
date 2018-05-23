package sidecar

import (
	"errors"
	"math"

	"dubbo-mesh/registry"
)

const (
	LB_Random  = iota
	LB_RR
	LB_WRR
	LB_LLatest
	LB_LAvg
	LB_LActive
	LB_DRR
)

func lb(elector int) Banlancer {
	switch elector {
	case LB_Random:
		return &Random{}
	case LB_RR:
		return &RoundRobin{}
	case LB_WRR:
		return &WeightRoundRobin{}
	case LB_LLatest:
		return &LeastLatest{}
	case LB_LAvg:
		return &LeastAVG{}
	case LB_DRR:
		return &DynamicWeightRoundRobin{}
	case LB_LActive:
		return &LeastActive{}
	default:
		panic(errors.New("unknown load balancer"))
	}
}

type Banlancer interface {
	Init(endpoint []*Endpoint)
	Elect(endpoints []*Endpoint) *Endpoint
}

type Endpoint struct {
	*registry.Endpoint
	Meter       *Meter
	Active      int32
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
	TotalCount uint64 `json:"total_count,omitempty"` // 处理的总数
	ErrorCount uint64 `json:"error_count,omitempty"` // 错误数
	Error      uint64 `json:"error,omitempty"`       // 大于0，最近n次错误，=0 最近一次没有错误
	Latest     uint64 `json:"latest,omitempty"`      // RTT
	Max        uint64 `json:"max,omitempty"`
	Min        uint64 `json:"min,omitempty"`
	Total      uint64 `json:"total,omitempty"`
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

func NewEndpoint(endpoint *registry.Endpoint) *Endpoint {
	return &Endpoint{
		Endpoint: endpoint,
		Meter: &Meter{
			Min: math.MaxUint64,
		},
	}
}
