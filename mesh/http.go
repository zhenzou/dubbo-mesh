package mesh

import (
	"net/http"
	"net/url"
	"time"
	"fmt"
	"context"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
	"dubbo-mesh/log"
	"dubbo-mesh/dubbo"
)

func NewHttpClient() Client {
	return &HttpClient{
		&http.Client{
			Transport: &http.Transport{
				Proxy: nil,
			},
		},
	}
}

type HttpClient struct {
	client *http.Client
}

func (this *HttpClient) Invoke(endpoint *registry.EndPoint, invocation *Invocation) ([]byte, error) {
	form := url.Values{}
	form.Add(ParamInterface, invocation.Interface)
	form.Add(ParamMethod, invocation.Method)
	form.Add(ParamParamType, invocation.ParamType)
	form.Add(ParamParam, invocation.Param)
	resp, err := this.client.PostForm("http://"+endpoint.String(), form)
	if err != nil {
		return nil, err
	} else {
		data, err := util.ReadResponse(resp)
		return data, err
	}
}

func NewHttpServer(port int, dubbo *dubbo.Client) Server {
	return &HttpServer{client: dubbo, port: port}
}

type HttpServer struct {
	*http.Server
	port   int
	client *dubbo.Client
}

func (this *HttpServer) Invocations() <-chan Invocation {
	// TODO
	return nil
}

func (this *HttpServer) Run() error {
	log.Info("server start to run on port ", this.port)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: http.HandlerFunc(this.invoke),
	}
	return this.ListenAndServe()
}

func (this *HttpServer) invoke(w http.ResponseWriter, req *http.Request) {
	interfaceName := req.FormValue(ParamInterface)
	method := req.FormValue(ParamMethod)
	paramType := req.FormValue(ParamParamType)
	param := req.FormValue(ParamParam)
	resp, err := this.client.Invoke(interfaceName, method, paramType, param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if resp.Error() != nil {
		log.Warn(resp.Error().Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		//w.Header().Set("Content-Type", "application/json; charset=utf-8")
		log.Debug("resp:", string(resp.Body()))
		w.Write(resp.Body())
	}
}

func (this *HttpServer) Shutdown() error {
	log.Info("server start to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := this.Server.Shutdown(ctx)
	return err
}
