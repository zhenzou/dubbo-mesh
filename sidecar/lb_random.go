package sidecar

import (
	"math/rand"
)

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
