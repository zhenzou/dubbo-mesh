package sidecar

import (
	"net/http"

	"dubbo-mesh/dubbo"
	"dubbo-mesh/log"
	"dubbo-mesh/derror"
)

func NewProvider(cfg *Config) *Provider {
	server := NewServer(cfg)
	client := dubbo.NewClient(cfg.DubboAddr)
	provider := &Provider{server, client}
	derror.Panic(provider.init())
	server.handler = provider.invoke
	return provider
}

type Provider struct {
	*Server
	client *dubbo.Client
}

func (this *Provider) init() error {
	return this.registry.Register(this.cfg.Service, this.cfg.ServerPort)
}

func (this *Provider) invoke(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue("interface")
	method := req.FormValue("method")
	paramType := req.FormValue("parameterTypesString")
	param := req.FormValue("parameter")
	resp, err := this.client.Invoke(interfaceName, method, paramType, param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if resp.Error() != nil {
		log.Warn(resp.Error().Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		log.Debug("resp:", string(resp.Body()))
		w.Write(resp.Body())
	}
}

func (this *Provider) Shutdown() {
	derror.Warn(this.Server.Shutdown())
	this.client.Shutdown()
}
