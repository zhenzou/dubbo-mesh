package registry

import (
	"fmt"
)

const (
	RootPath = "dubbo-mesh"
)

// 系统配置相关信息
type System struct {
	CpuNum      int
	TotalMemory int
	UsedMemory  int
}

type Endpoint struct {
	Host   string  `json:"host"`
	Port   int     `json:"-"`
	System *System `json:"-"`
}

func (this *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
}

type Registry interface {
	Register(serviceName string, port int) error
	Find(serviceName string) ([]*Endpoint, error)
}
