package sidecar

import (
	"net/http"
	"fmt"
	"time"
	"context"

	"dubbo-mesh/log"
)

type Config struct {
	Service    string
	DubboAddr  string
	ServerPort int
	Etcd       string
	Elector    int
}

// 封装http.Server
type Server struct {
	*http.Server
	port    int
	handler http.HandlerFunc
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (this *Server) Run() error {
	log.Info("server start to run on port ", this.port)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
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
