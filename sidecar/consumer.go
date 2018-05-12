package sidecar

import (
	"net/http"

	"dubbo-mesh/registry"
	"dubbo-mesh/derror"
	"dubbo-mesh/mesh"
	"dubbo-mesh/log"
	"dubbo-mesh/util"
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
		Elector:  elector(cfg.Elector),
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
	endpoints []*Endpoint
	Elector   Elector
}

func (this *Consumer) init() error {
	endpoints, err := this.registry.Find(this.cfg.Service)
	if err != nil {
		return err
	}
	log.Info("get service:", util.ToJsonStr(endpoints))
	this.endpoints = make([]*Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		this.endpoints[i] = NewEndpoint(endpoint)
	}
	this.Elector.Init(this.endpoints)
	return nil
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
	if len(this.endpoints) == 1 {
		return this.endpoints[0].Endpoint
	}
	return this.Elector.Elect(this.endpoints).Endpoint
}
