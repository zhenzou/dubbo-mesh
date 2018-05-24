// +build prod

package sidecar

import (
	"dubbo-mesh/mesh"
)

func (this *Consumer) asyncRecord() {
}

func (this *Consumer) invoke(inv *mesh.Invocation) ([]byte, error) {
	endpoint := this.Elect()
	data, err := this.Invoke(endpoint.Endpoint, inv)
	return data, err
}
