package sidecar

import (
	"sync/atomic"
	"net/http"
	"sync"

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
		Elector:  elector(cfg.Balancer),
		Client:   mesh.NewTcpClient(),
		rtts:     make(chan *Rtt, 200),
	}
	derror.Panic(consumer.init())
	server.handler = consumer.invoke
	return consumer
}

var (
	invPool = sync.Pool{
		New: func() interface{} {
			return &mesh.Invocation{}
		},
	}
)

type Consumer struct {
	mesh.Client
	*Server
	cfg       *Config
	registry  registry.Registry
	endpoints []*Endpoint
	Elector   Banlancer
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
	//go this.asyncRecord()
	return nil
}

func (this *Consumer) invoke(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue("interface")
	method := req.FormValue("method")
	paramType := req.FormValue("parameterTypesString")
	param := req.FormValue("parameter")
	inv := invPool.Get().(*mesh.Invocation)
	defer invPool.Put(inv)
	inv.Interface = interfaceName
	inv.Method = method
	inv.ParamType = paramType
	inv.Param = param
	// TODO retry,会影响性能
	endpoint := this.Elect()
	//log.Debug("status:", util.ToJsonStr(endpoint.Meter))
	//start := time.Now()dock
	data, err := this.Invoke(endpoint.Endpoint, inv)
	//end := time.Now()
	//this.rtts <- &Rtt{Endpoint: endpoint, Rtt: end.Sub(start).Nanoseconds(), Error: err}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

func (this *Consumer) syncRecord(endpoint *Endpoint, nano uint64, err error) {
	atomic.AddUint64(&endpoint.Meter.TotalCount, 1)
	atomic.AddUint64(&endpoint.Meter.Total, nano)
	atomic.StoreUint64(&endpoint.Meter.Latest, nano)
	if nano < endpoint.Meter.Min {
		atomic.StoreUint64(&endpoint.Meter.Min, nano)
	}
	if nano > endpoint.Meter.Max {
		atomic.StoreUint64(&endpoint.Meter.Max, nano)
	}
	if err != nil {
		atomic.AddUint64(&endpoint.Meter.ErrorCount, 1)
		atomic.AddUint64(&endpoint.Meter.Error, 1)
	} else {
		atomic.StoreUint64(&endpoint.Meter.Error, 0)
	}
}

func (this *Consumer) asyncRecord() {
	for rtt := range this.rtts {
		endpoint := rtt.Endpoint
		nano := uint64(rtt.Rtt)
		err := rtt.Error
		endpoint.Meter.TotalCount += 1
		endpoint.Meter.Total += nano
		endpoint.Meter.Latest = nano
		if nano < endpoint.Meter.Min {
			endpoint.Meter.Min = nano
		}
		if nano > endpoint.Meter.Max {
			endpoint.Meter.Max = nano
		}
		if err != nil {
			endpoint.Meter.ErrorCount += 1
			endpoint.Meter.Error += 1
		} else {
			endpoint.Meter.Error = 0
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
