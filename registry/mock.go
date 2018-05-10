package registry

func NewMock(endpoints []*Endpoint) Registry {
	return &Mock{endpoints}
}

func DefaultMock() Registry {
	return NewMock([]*Endpoint{
		{
			Host: "127.0.0.1",
			Port: 30001,
		},
	})
}

// mock 单机测试
type Mock struct {
	endpoints []*Endpoint
}

func (this *Mock) Register(serviceName string, port int) error {
	return nil
}

func (this *Mock) Find(serviceName string) ([]*Endpoint, error) {
	return this.endpoints, nil
}
