package dubbo

import (
	"net/http"
	"fmt"
	"time"
	"context"

	"dubbo-mesh/log"
)

type Config struct {
	DubboAddr  string
	ServerPort int
	Etcd       string
}

var (
	client *Client
)

type Server struct {
	*http.Server
	handler http.HandlerFunc
}

func NewServer(cfg *Config) *Server {
	client = NewClient(cfg.DubboAddr)

	return &Server{}
}

func (this *Server) Run(port int) error {
	log.Info("server start to run on port ", port)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(invoke),
	}
	return this.ListenAndServe()
}

func (this *Server) Shutdown() error {
	log.Info("server start to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return this.Server.Shutdown(ctx)
}

func invoke(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue("interface")
	method := req.FormValue("method")
	paramType := req.FormValue("parameterTypesString")
	param := req.FormValue("parameter")
	log.Debugf("%s/%s/%s/%s", interfaceName, method, paramType, param)
	resp, err := client.Invoke(interfaceName, method, paramType, param)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resp.Payload)
	}
}
