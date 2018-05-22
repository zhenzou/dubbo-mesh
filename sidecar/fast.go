package sidecar

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"dubbo-mesh/log"
)

// 封装http.Server
type FastServer struct {
	port    int
	handler fasthttp.RequestHandler
}

func NewFastServer(port int) Server {
	return &FastServer{port: port}
}

func (this *FastServer) Run() error {
	log.Info("server start to run on port ", this.port)
	return fasthttp.ListenAndServe(fmt.Sprintf(":%d", this.port), this.handler)
}

func (this *FastServer) Shutdown() error {
	log.Info("server start to shutdown")
	return nil
}
