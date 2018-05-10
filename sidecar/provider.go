package sidecar

import (
	"dubbo-mesh/dubbo"
	"dubbo-mesh/derror"
	"dubbo-mesh/registry"
	"dubbo-mesh/mesh"
)

func NewMockProvider(cfg *Config) *Provider {
	client := dubbo.NewClient(cfg.DubboAddr)
	server := mesh.NewHttpServer(cfg.ServerPort, client)
	provider := &Provider{server, cfg, registry.DefaultMock()}
	derror.Panic(provider.init())
	return provider
}

func NewProvider(cfg *Config) *Provider {
	client := dubbo.NewClient(cfg.DubboAddr)
	server := mesh.NewHttpServer(cfg.ServerPort, client)
	provider := &Provider{server, cfg, registry.NewEtcdFromAddr(cfg.Etcd)}
	derror.Panic(provider.init())
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
