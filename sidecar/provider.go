package sidecar

import (
	"dubbo-mesh/dubbo"
	"dubbo-mesh/derror"
	"dubbo-mesh/registry"
	"dubbo-mesh/mesh"
	"dubbo-mesh/log"
)

func NewMockProvider(cfg *Config) *Provider {
	return newProvider(cfg, registry.DefaultMock())
}

func NewProvider(cfg *Config) *Provider {
	return newProvider(cfg, registry.NewEtcdFromAddr(cfg.Etcd))
}

func newProvider(cfg *Config, registry registry.Registry) *Provider {
	client := dubbo.NewClient(cfg.DubboAddr)
	log.Debug("server start")
	server := mesh.NewTcpServer(cfg.ServerPort, client)
	log.Debug("server new ")
	provider := &Provider{server, cfg, registry}
	derror.Panic(provider.init())
	log.Debug("server init")
	log.Debug("server stop")
	return provider
}

type Provider struct {
	mesh.Server
	cfg      *Config
	registry registry.Registry
}

func (this *Provider) init() error {
	return this.registry.Register(this.cfg.Service, this.cfg.ServerPort)
}
