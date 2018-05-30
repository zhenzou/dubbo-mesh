// +build !prod

package sidecar

import (
	"time"
	"sync/atomic"

	"dubbo-mesh/mesh"
)

func (this *Consumer) invoke(inv *mesh.Invocation) ([]byte, error) {
	// TODO retry.会影响性能
	endpoint := this.Elect()
	atomic.AddInt32(&endpoint.Active, 1)
	start := time.Now()
	data, err := this.Invoke(endpoint.Endpoint, inv)
	atomic.AddInt32(&endpoint.Active, -1)
	end := time.Now()
	this.syncRecord(endpoint, uint64(end.Sub(start).Nanoseconds()/1000000))
	return data, err
}
