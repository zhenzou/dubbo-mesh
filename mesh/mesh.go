package mesh

import (
	"sync"

	"dubbo-mesh/registry"
	"dubbo-mesh/json"
)

type Protocol int

const (
	Tcp  Protocol = iota + 1
	Http
	Kcp
)

var (
	invs = sync.Pool{
		New: func() interface{} {
			return &Invocation{}
		},
	}
)

func NewInv() *Invocation {
	return invs.Get().(*Invocation)
}

func ReleaseInv(inv *Invocation) {
	invs.Put(inv)
}

//  暂时只考虑Dubbo协议
type Invocation struct {
	Interface string `json:"i"`
	Method    string `json:"m"`
	ParamType string `json:"pt"`
	Param     string `json:"p"`
}

func (this *Invocation) Data() []byte {
	data, _ := json.Marshal(this)
	return data
}

type Client interface {
	Invoke(endpoint *registry.Endpoint, invocation *Invocation) ([]byte, error)
}

type Server interface {
	Run() error
	Shutdown() error
}
