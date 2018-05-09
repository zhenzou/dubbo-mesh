package registry

func NewMock(endpoints []*EndPoint) Registry {
	return &Mock{endpoints}
}

func DefaultMock() Registry {
	return NewMock([]*EndPoint{
		{
			Host: "127.0.0.1",
			Port: 30001,
		},
	})
}

// mock 单机测试
type Mock struct {
	endpoints []*EndPoint
}

func (this *Mock) Register(serviceName string, port int) error {
	return nil
}

func (this *Mock) Find(serviceName string) ([]*EndPoint, error) {
	return this.endpoints, nil
}
