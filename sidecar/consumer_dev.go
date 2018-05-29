// +build !prod

package sidecar

import (
	"time"
	"sync/atomic"

	"dubbo-mesh/mesh"
)

func (this *Consumer) asyncRecord() {
	count := 0
	for rtt := range this.rtts {
		count++
		endpoint := rtt.Endpoint
		mill := uint64(rtt.Rtt)
		endpoint.Meter.Count += 1
		endpoint.Meter.Total += mill
		endpoint.Meter.Latest = mill
		if mill < endpoint.Meter.Min {
			endpoint.Meter.Min = mill
		}
		if mill > endpoint.Meter.Max {
			endpoint.Meter.Max = mill
		}
		if count > 0 && count%10000 == 0 {
			this.print()
		}
	}
}

func (this *Consumer) invoke(inv *mesh.Invocation) ([]byte, error) {
	// TODO retry.会影响性能
	endpoint := this.Elect()
	atomic.AddInt32(&endpoint.Active, 1)
	start := time.Now()
	data, err := this.Invoke(endpoint.Endpoint, inv)
	atomic.AddInt32(&endpoint.Active, -1)
	end := time.Now()
	this.rtts <- &Rtt{Endpoint: endpoint, Rtt: end.Sub(start).Nanoseconds() / 1000000}
	return data, err
}
