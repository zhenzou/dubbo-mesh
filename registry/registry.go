package registry

const (
	RootPath = "dubbo-mesh"
)

// 系统配置相关信息
type System struct {
	CpuNum int
	Memory int
	Name   string `json:"name"`
}

type Endpoint struct {
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	System *System `json:"system"`
}

type Registry interface {
	Register(serviceName string, port int) error
	Find(serviceName string) ([]*Endpoint, error)
}
