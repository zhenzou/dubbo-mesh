package sidecar

import (
	"net/http"

	"dubbo-mesh/registry"
	"dubbo-mesh/derror"
	"dubbo-mesh/mesh"
	"dubbo-mesh/log"
	"dubbo-mesh/util"
	"time"
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
		rtts:     make(chan *Rtt, 200),
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
	rtts      chan *Rtt
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
	go this.record()
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
	// TODO retry,会影响性能
	endpoint := this.Elect()
	log.Debug("status:", util.ToJsonStr(endpoint.Status))
	start := time.Now()
	data, err := this.Invoke(endpoint.Endpoint, inv)
	end := time.Now()
	this.rtts <- &Rtt{Endpoint: endpoint, Rtt: end.Sub(start).Nanoseconds(), Error: err}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

// TODO 异步，串行化
func (this *Consumer) record() {
	for rtt := range this.rtts {
		endpoint := rtt.Endpoint
		nano := uint64(rtt.Rtt)
		err := rtt.Error
		endpoint.Status.Count += 1
		endpoint.Status.Total += nano
		endpoint.Status.Latest = nano
		if nano < endpoint.Status.Min {
			endpoint.Status.Min = nano
		}
		if nano > endpoint.Status.Max {
			endpoint.Status.Max = nano
		}
		if err != nil {
			endpoint.Status.ErrorCount += 1
			endpoint.Status.Error += 1
		} else {
			endpoint.Status.Error = 0
		}
	}

}

// 负载均衡，选择其中一个
// TODO 优化策略
func (this *Consumer) Elect() *Endpoint {
	if len(this.endpoints) == 1 {
		return this.endpoints[0]
	}
	return this.Elector.Elect(this.endpoints)
}
