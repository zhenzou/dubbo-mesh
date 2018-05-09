package sidecar

import (
	"net/http"
	"fmt"
	"time"
	"context"

	"dubbo-mesh/log"
	"dubbo-mesh/registry"
)

type Config struct {
	Service    string
	DubboAddr  string
	ServerPort int
	Etcd       string
}

type Server struct {
	*http.Server
	cfg      *Config
	registry registry.Registry
	handler  http.HandlerFunc
}

func NewServerWithRegistry(cfg *Config, registry registry.Registry) *Server {
	return &Server{cfg: cfg, registry: registry}
}

func NewServer(cfg *Config) *Server {
	return &Server{cfg: cfg, registry: registry.NewEtcdFromAddr(cfg.Etcd)}
}

func (this *Server) Run() error {
	log.Info("server start to run on port ", this.cfg.ServerPort)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", this.cfg.ServerPort),
		Handler: this.handler,
	}
	return this.ListenAndServe()
}

func (this *Server) Shutdown() error {
	log.Info("server start to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := this.Server.Shutdown(ctx)
	return err
}
