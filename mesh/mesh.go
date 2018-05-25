package mesh

import (
	"bytes"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
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

func (this *Invocation) Data() []byte {
	data := bytes.Join([][]byte{util.StringToBytes(this.Interface), util.StringToBytes(this.Method),
		util.StringToBytes(this.ParamType), util.StringToBytes(this.Param)}, []byte("\n"))
	return data
}

type Client interface {
	Invoke(endpoint *registry.Endpoint, invocation *Invocation) ([]byte, error)
}

type Server interface {
	Run() error
	Shutdown() error
}
