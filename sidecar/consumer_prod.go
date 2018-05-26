// +build prod

package sidecar

import (
	"sync/atomic"

	"dubbo-mesh/mesh"
)

func (this *Consumer) asyncRecord() {
}

func (this *Consumer) invoke(inv *mesh.Invocation) ([]byte, error) {
	endpoint := this.Elect()
	atomic.AddInt32(&endpoint.Active, 1)
	data, err := this.Invoke(endpoint.Endpoint, inv)
	atomic.AddInt32(&endpoint.Active, -1)
	return data, err
}
