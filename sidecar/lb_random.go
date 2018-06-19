package sidecar

import (
	"math/rand"
)

type Random struct {
	baseBalancer
	total int
}

func (this *Random) Init(endpoints []*Endpoint) {
	this.baseBalancer.Init(endpoints)
	this.total = len(endpoints)
}

func (this *Random) Elect() *Endpoint {
	i := rand.Intn(this.total)
	return this.endpoints[i]
}

// 加权随机,没有实现add
type WeightRandom struct {
	baseBalancer
	indexes map[int]*Endpoint
	total   int
}

func (this *WeightRandom) Init(endpoints []*Endpoint) {
	weights := make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weights[endpoint] = this.weight(endpoint)
	}
	gcd := this.gcd(weights)
	for k, originWeight := range weights {
		weight := originWeight / gcd
		weights[k] = weight
		this.total += weight
	}
	index := 0
	indexes := make(map[int]*Endpoint, this.total)
	for end, weight := range weights {
		for i := 0; i < weight; i++ {
			indexes[index] = end
			index++
		}
	}
	this.indexes = indexes
}

func (this *WeightRandom) Elect() *Endpoint {
	n := rand.Intn(this.total)
	return this.indexes[n]
}
