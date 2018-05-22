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
	Balancer   int
	Server     int
}

type Server interface {
	Run() error
	Shutdown() error
}

// 封装http.Server
type HttpServer struct {
	*http.Server
	port    int
	handler http.HandlerFunc
}

func NewServer(port int) Server {
	return &HttpServer{port: port}
}

func (this *HttpServer) Run() error {
	log.Info("server start to run on port ", this.port)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: this.handler,
	}
	return this.ListenAndServe()
}

func (this *HttpServer) Shutdown() error {
	log.Info("server start to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := this.Server.Shutdown(ctx)
	return err
}
