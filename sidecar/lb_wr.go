package sidecar

import (
	"math/rand"

	"dubbo-mesh/util"
)

// 加权随机
type WeightRandom struct {
	indexes map[int]*Endpoint
	total   int
}

func (this *WeightRandom) Init(endpoints []*Endpoint) {
	weights := make(map[*Endpoint]int)
	for _, endpoint := range endpoints {
		weights[endpoint] = this.calculateWrr(endpoint)
	}
	gcd := this.weightGcd(weights)
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

func (this *WeightRandom) Elect(endpoints []*Endpoint) *Endpoint {
	n := rand.Intn(this.total)
	return this.indexes[n]
}

func (this *WeightRandom) weightGcd(weights map[*Endpoint]int) int {
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
func (this *WeightRandom) calculateWrr(status *Endpoint) int {
	return status.System.Memory
}
