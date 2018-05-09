package registry

import (
	"fmt"
)

const (
	RootPath = "dubbo-mesh"
)

type Status struct {
	System *System
	Alive  int // 活跃连接数
	Rate   int // 处理速率
	Rtt    *Rtt
}

// 系统配置相关信息
type System struct {
	Core        int
	TotalMemory int
	FreeMemory  int
}

// RTT统计
type Rtt struct {
	Max int
	Min int
	Avg int
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
