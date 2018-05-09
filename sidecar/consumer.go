package sidecar

import (
	"net/http"
	"math/rand"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
	"dubbo-mesh/derror"
)

func NewMockConsumer(cfg *Config) *Consumer {
	server := NewServerWithRegistry(cfg, registry.DefaultMock())
	consumer := &Consumer{Server: server}
	derror.Panic(consumer.init())
	server.handler = consumer.invoke
	return consumer
}

func NewConsumer(cfg *Config) *Consumer {
	server := NewServer(cfg)
	consumer := &Consumer{Server: server}
	derror.Panic(consumer.init())
	server.handler = consumer.invoke
	return consumer
}

type Consumer struct {
	*Server
	endpoints []*registry.EndPoint
}

func (this *Consumer) init() error {
	var err error
	this.endpoints, err = this.registry.Find(this.cfg.Service)
	return err
}

func (this *Consumer) invoke(w http.ResponseWriter, req *http.Request) {
	endPoint := this.Elect()
	req.ParseForm()
	resp, err := http.PostForm("http://"+endPoint.String(), req.Form)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		data, _ := util.ReadResponse(resp)
		w.Write(data)
	}
}

// 负载均衡，选择其中一个
// TODO 优化策略
func (this *Consumer) Elect() *registry.EndPoint {
	i := rand.Intn(len(this.endpoints))
	return this.endpoints[i]
}
