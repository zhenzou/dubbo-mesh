package mesh

import (
	"dubbo-mesh/registry"
)

type Protocol int

const (
	Tcp  Protocol = iota + 1
	Http
	Kcp
)

//  暂时只考虑Dubbo协议
type Invocation struct {
	Interface string `json:"i"`
	Method    string `json:"m"`
	ParamType string `json:"pt"`
	Param     string `json:"p"`
}

type Client interface {
	Invoke(endpoint *registry.Endpoint, invocation *Invocation) ([]byte, error)
}

type Server interface {
	Invocations() <-chan Invocation
	Run() error
	Shutdown() error
}
