package registry

import (
	"fmt"
)

const (
	RootPath = "dubbo-mesh"
)

type Status struct {
	Alive  int
	Core   int
	Memory int
	Rate   int
}

type EndPoint struct {
	Host  string
	Port  int
	Value *Status
}

func (this *EndPoint) String() string {
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
}

type Registry interface {
	Register(serviceName string, port int) error
	Find(serviceName string) ([]*EndPoint, error)
}
