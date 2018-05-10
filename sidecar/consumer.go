package sidecar

import (
	"net/http"
	"math/rand"

	"dubbo-mesh/registry"
	"dubbo-mesh/derror"
	"dubbo-mesh/mesh"
)

func NewMockConsumer(cfg *Config) *Consumer {
	consumer := newConsumer(cfg, registry.DefaultMock())
	return consumer
}

func NewConsumer(cfg *Config) *Consumer {
	consumer := newConsumer(cfg, registry.NewEtcdFromAddr(cfg.Etcd))
	return consumer
}

func newConsumer(cfg *Config, registry registry.Registry) *Consumer {
	server := NewServer(cfg.ServerPort)
	consumer := &Consumer{
		cfg:      cfg,
		Server:   server,
		registry: registry,
		Client:   mesh.NewTcpClient(),
	}
	derror.Panic(consumer.init())
	server.handler = consumer.invoke
	return consumer
}

type Consumer struct {
	mesh.Client
	*Server
	cfg       *Config
	registry  registry.Registry
	endpoints []*registry.Endpoint
}

func (this *Consumer) init() error {
	var err error
	this.endpoints, err = this.registry.Find(this.cfg.Service)
	return err
}

func (this *Consumer) invoke(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue("interface")
	method := req.FormValue("method")
	paramType := req.FormValue("parameterTypesString")
	param := req.FormValue("parameter")
	inv := &mesh.Invocation{
		Interface: interfaceName,
		Method:    method,
		ParamType: paramType,
		Param:     param,
	}
	endpoint := this.Elect()
	data, err := this.Invoke(endpoint, inv)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

// 负载均衡，选择其中一个
// TODO 优化策略
func (this *Consumer) Elect() *registry.Endpoint {
	i := rand.Intn(len(this.endpoints))
	return this.endpoints[i]
}
