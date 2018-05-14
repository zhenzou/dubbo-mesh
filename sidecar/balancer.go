package sidecar

import (
	"math/rand"
	"sync/atomic"
	"errors"
	"math"
	"time"

	"dubbo-mesh/registry"
	"dubbo-mesh/log"
	"dubbo-mesh/util"
)

const (
	ElectorRandom = iota
	ElectorRR
	ElectorWRR
)

func elector(elector int) Banlancer {
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

type Banlancer interface {
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

// TODO 解决动态变化
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

type WrrRandom struct {
	weights map[*Endpoint]int
	index   int
	max     int
	cw      int // 当前权重
}

func (this *WrrRandom) Init(endpoints []*Endpoint) {
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

func (this *WrrRandom) Elect(endpoints []*Endpoint) *Endpoint {

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
	return nil
}

func (r *WrrRandom) weightGcd() int {
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
func (this *WrrRandom) calculateWrr(status *Endpoint) int {
	return status.System.TotalMemory
}

// 动态权重变化
type DrrRandom struct {
	weights map[*Endpoint]*int
	total   int
}

func (this *DrrRandom) Init(endpoints []*Endpoint) {
	this.weights = make(map[*Endpoint]*int)
	for _, endpoint := range endpoints {
		weight := this.initCalculateWrr(endpoint)
		this.weights[endpoint] = &weight
		this.total += this.initCalculateWrr(endpoint)
	}
	go this.cronCalculateDrr()
}

func (this *DrrRandom) Elect(endpoints []*Endpoint) *Endpoint {
	w := rand.Intn(this.total)
	for endpoint, weight := range this.weights {
		w -= *weight
		if w < 0 {
			return endpoint
		}
	}
	return nil
}

// 简单的计算权重，只考虑系统配置
func (this *DrrRandom) initCalculateWrr(endpoint *Endpoint) int {
	return endpoint.System.CpuNum + endpoint.System.TotalMemory/100000
}

// 动态的计算权重，考虑系统配置和运行状态
func (this *DrrRandom) cronCalculateDrr() {
	tick := time.Tick(time.Second)
	for _ = range tick {
		for endpoint := range this.weights {
			init := endpoint.System.CpuNum + endpoint.System.TotalMemory/100000
			dw := this.dw(endpoint.Status)
			*this.weights[endpoint] = (init*30 + dw*70) / 100
		}
	}
}

// 运行状态权重，归一到100
func (this *DrrRandom) dw(status *Status) int {
	weight := 5*status.Latest + 2*status.Min - 2*status.Max + 1*status.Avg() - 1000000*status.ErrorCount
	if status.Error > 0 {
		weight -= 100000000
	}
	return int(weight % 100)
}

func NewEndpoint(endpoint *registry.Endpoint) *Endpoint {
	return &Endpoint{
		Endpoint: endpoint,
		Status: &Status{
			Min: math.MaxUint64,
		},
	}
}

type Endpoint struct {
	*registry.Endpoint
	Status *Status
}

type Rtt struct {
	Endpoint *Endpoint
	Rtt      int64
	Error    error
}

type Status struct {
	Count      uint64 // 处理的总数
	ErrorCount uint64 // 错误数
	Error      uint64 // 大于0，最近n次错误，=0 最近一次没有错误
	Latest     uint64 // RTT
	Max        uint64
	Min        uint64
	Total      uint64
}

func (this *Status) Avg() uint64 {
	return this.Total / this.Count
}
