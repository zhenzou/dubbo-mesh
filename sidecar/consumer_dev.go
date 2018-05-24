// +build !prod

package sidecar

import (
	"dubbo-mesh/log"
)

func (this *Consumer) asyncRecord() {
	count := 0
	for rtt := range this.rtts {
		count++
		endpoint := rtt.Endpoint
		nano := uint64(rtt.Rtt)
		err := rtt.Error
		endpoint.Meter.TotalCount += 1
		endpoint.Meter.Total += nano
		endpoint.Meter.Latest = nano
		if nano < endpoint.Meter.Min {
			endpoint.Meter.Min = nano
		}
		if nano > endpoint.Meter.Max {
			endpoint.Meter.Max = nano
		}
		if err != nil {
			endpoint.Meter.ErrorCount += 1
			endpoint.Meter.Error += 1
		} else {
			endpoint.Meter.Error = 0
		}
		if count > 0 && count%10000 == 0 {
			log.Info(rtt.Endpoint.String())
		}
	}
}
