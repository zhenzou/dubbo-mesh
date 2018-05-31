package sidecar

import (
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"

	"dubbo-mesh/derror"
	"dubbo-mesh/log"
	"dubbo-mesh/mesh"
	"dubbo-mesh/registry"
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
	var server Server
	if cfg.Server == 0 {
		server = NewServer(cfg.ServerPort)
	} else {
		server = NewFastServer(cfg.ServerPort)
	}
	consumer := &Consumer{
		cfg:      cfg,
		Server:   server,
		registry: registry,
		Balancer: lb(cfg.Balancer),
		Client:   mesh.NewTcpClient(),
		rtts:     make(chan *Rtt, 256),
	}
	derror.Panic(consumer.init())
	switch s := server.(type) {
	case *HttpServer:
		s.handler = consumer.httpHandler
	case *FastServer:
		s.handler = consumer.fastHandler
	default:
		panic(errors.New("wrong server"))
	}
	return consumer
}

type Consumer struct {
	Server
	mesh.Client
	cfg       *Config
	registry  registry.Registry
	endpoints []*Endpoint
	Balancer  Banlancer
	rtts      chan *Rtt
}

func (this *Consumer) init() error {
	endpoints, err := this.registry.Find(this.cfg.Service)
	if err != nil {
		return err
	}
	log.Info("providers:", util.ToJsonStr(endpoints))
	this.endpoints = make([]*Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		this.endpoints[i] = NewEndpoint(endpoint)
	}
	this.Balancer.Init(this.endpoints)
	return nil
}

func (this *Consumer) httpHandler(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue("interface")
	method := req.FormValue("method")
	paramType := req.FormValue("parameterTypesString")
	param := req.FormValue("parameter")

	inv := mesh.NewInv()
	defer mesh.ReleaseInv(inv)

	inv.Interface = interfaceName
	inv.Method = method
	inv.ParamType = paramType
	inv.Param = param

	data, err := this.invoke(inv)

	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

func (this *Consumer) fastHandler(ctx *fasthttp.RequestCtx) {
	interfaceName := ctx.FormValue("interface")
	method := ctx.FormValue("method")
	paramType := ctx.FormValue("parameterTypesString")
	param := ctx.FormValue("parameter")

	inv := mesh.NewInv()
	defer mesh.ReleaseInv(inv)

	inv.Interface = util.BytesToString(interfaceName)
	inv.Method = util.BytesToString(method)
	inv.ParamType = util.BytesToString(paramType)
	inv.Param = util.BytesToString(param)

	data, err := this.invoke(inv)

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
	} else {
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(data)
	}
}

func (this *Consumer) invoke(inv *mesh.Invocation) ([]byte, error) {
	// TODO retry.会影响性能
	endpoint := this.Elect()
	atomic.AddInt32(&endpoint.Meter.Active, 1)
	start := time.Now()
	data, err := this.Invoke(endpoint.Endpoint, inv)
	atomic.AddInt32(&endpoint.Meter.Active, -1)
	end := time.Now()
	this.syncRecord(endpoint, uint64(end.Sub(start).Nanoseconds()/1000000))
	return data, err
}

func (this *Consumer) syncRecord(endpoint *Endpoint, mill uint64) {
	atomic.AddUint64(&endpoint.Meter.Count, 1)
	atomic.AddUint64(&endpoint.Meter.Total, mill)
	atomic.StoreUint64(&endpoint.Meter.Latest, mill)
	if mill < endpoint.Meter.Min {
		atomic.StoreUint64(&endpoint.Meter.Min, mill)
	}
	if mill > endpoint.Meter.Max {
		atomic.StoreUint64(&endpoint.Meter.Max, mill)
	}
}

func (this *Consumer) print() {
	for _, endpoint := range this.endpoints {
		log.Info(endpoint.String())
	}
}

func (this *Consumer) Shutdown() error {
	this.print()
	return this.Server.Shutdown()
}

func (this *Consumer) Elect() *Endpoint {
	if len(this.endpoints) == 1 {
		return this.endpoints[0]
	}
	return this.Balancer.Elect(this.endpoints)
}
