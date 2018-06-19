package sidecar

import (
	"errors"
	"sync"

	"gopkg.in/eapache/queue.v1"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
)

const (
	LB_Random = iota + 1
	LB_RR
	LB_WRR
	LB_LAvg
	LB_LA
	LB_LLA
	LB_WLLA
)



func lb(elector int) Balancer {
	switch elector {
	case LB_Random:
		return &Random{}
	case LB_RR:
		return &RoundRobin{}
	case LB_WRR:
		return &WeightRoundRobin{}
	case LB_LAvg:
		return &LeastAVG{}
	case LB_LA:
		return &LeastActive{}
	case LB_LLA:
		return &LeastLatestAvg{}
	case LB_WLLA:
		return &WeightLeastLatestAvg{}
	default:
		panic(errors.New("unknown load balancer"))
	}
}

type Balancer interface {
	Init(endpoints []*Endpoint)
	Record(endpoint *Endpoint, latency uint64)
	Add(endpoint *Endpoint)
	Count() int
	Elect() *Endpoint
}

type baseBalancer struct {
	endpoints []*Endpoint
}

func (this *baseBalancer) Record(endpoint *Endpoint, latency uint64) {
	meter := endpoint.Meter
	meter.Mtx.Lock()
	meter.Count += 1
	meter.Latest = latency
	meter.Total += latency
	if latency < meter.Min {
		meter.Min = latency
	} else if latency > meter.Max {
		meter.Max = latency
	}
	meter.Mtx.Unlock()
}

func (this *baseBalancer) Init(endpoints []*Endpoint) {
	this.endpoints = endpoints
}

func (this *baseBalancer) Add(endpoint *Endpoint) {
	this.endpoints = append(this.endpoints, endpoint)
}

func (this *baseBalancer) Count() int {
	return len(this.endpoints)
}

func (this *baseBalancer) Elect() *Endpoint {
	// impl
	return nil
}

func (this *baseBalancer) gcd(weights map[*Endpoint]int) int {
	divisor := -1
	for _, s := range weights {
		if divisor == -1 {
			divisor = s
		} else {
			divisor = util.Gcd(divisor, s)
		}
	}
	return divisor
}

// 简单的计算权重，暂时 就把内存做为权重
func (this *baseBalancer) weight(status *Endpoint) int {
	return status.System.Memory
}

type Endpoint struct {
	*registry.Endpoint
	Meter *Meter
}

func (this *Endpoint) String() string {
	m := make(map[string]interface{}, 4)
	m["name"] = this.System.Name
	m["avg"] = this.Meter.Avg()
	m["meter"] = this.Meter
	return util.ToJsonStr(m)
}

type Meter struct {
	Mtx        sync.Mutex
	Queue      *queue.Queue
	Latest     uint64 `json:"latest"`
	Min        uint64 `json:"min"`
	Max        uint64 `json:"max"`
	Error      bool   `json:"error"`
	ErrorCount uint64 `json:"error_count"`
	Count      uint64 `json:"count,omitempty"`  // 已处理的总数
	Active     int32  `json:"active,omitempty"` // 当前连接数
	Total      uint64 `json:"total,omitempty"`
}

// 平均值
func (this *Meter) Avg() uint64 {
	if this.Count == 0 {
		return 0
	}
	return this.Total / this.Count
}

func NewEndpoint(endpoint *registry.Endpoint) *Endpoint {
	return &Endpoint{
		Endpoint: endpoint,
		Meter: &Meter{
			Queue: queue.New(),
		},
	}
}
