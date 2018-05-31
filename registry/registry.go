package registry

import (
	"fmt"
	"strconv"
	"strings"
)

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
	addr   string  `json:"-"`
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	System *System `json:"system"`
}

func NewEndpoint(addr string) (*Endpoint, error) {
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("wrong addr %s", addr)
	}
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("wrong addr %s", addr)
	}
	return &Endpoint{Host: split[0], Port: port}, nil
}

func (this *Endpoint) Addr() string {
	return this.addr
}

type Registry interface {
	Register(serviceName string, port int) error
	Find(serviceName string) ([]*Endpoint, error)
}
