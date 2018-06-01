package sidecar

import (
	"errors"
	"math"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
)

const (
	LB_Random = iota
	LB_RR
	LB_WRR
	LB_LLatest
	LB_LAvg
	LB_LActive
	LB_WLActive
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
	case LB_LActive:
		return &LeastActive{}
	case LB_WLActive:
		return &WeightLeastActive{}
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
	Meter *Meter
}

func (this *Endpoint) String() string {
	m := make(map[string]interface{}, 3)
	m["name"] = this.System.Name
	m["avg"] = this.Meter.Avg()
	m["meter"] = this.Meter
	return util.ToJsonStr(m)
}

type Meter struct {
	Count  uint64 `json:"count,omitempty"`  // 已处理的总数
	Active int32  `json:"active,omitempty"` // 当前连接数
	Latest uint64 `json:"latest,omitempty"` // RTT
	Max    uint64 `json:"max,omitempty"`
	Min    uint64 `json:"min,omitempty"`
	Total  uint64 `json:"total,omitempty"`
}

// 平均RTT
func (this *Meter) Avg() uint64 {
	if this.Count == 0 {
		return this.Total
	}
	return this.Total / this.Count
}

func NewEndpoint(endpoint *registry.Endpoint) *Endpoint {
	return &Endpoint{
		Endpoint: endpoint,
		Meter: &Meter{
			Min: math.MaxUint64,
		},
	}
}
