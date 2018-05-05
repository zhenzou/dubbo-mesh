package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"dubbo-mesh/log"
)

const (
	HttpMethodAny = "ANY"
)

type Router struct {
	Method   string
	Path     string
	Handler  gin.HandlerFunc
	Children []*Router
}

func (this *Router) AppendChild(method, path string, handler gin.HandlerFunc) {
	child := Router{Path: path, Handler: handler, Method: method}
	this.Children = append(this.Children, &child)
}

// 封装gin，增加简便的添加Rest方法
type Server struct {
	*gin.Engine
	*http.Server
}

func NewServer() *Server {
	g := NewGin()
	return &Server{Engine: g}
}

func (this *Server) Run(port int) error {
	log.Info("server start to run on port ", port)
	this.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: this.Engine,
	}
	return this.ListenAndServe()
}

func (this *Server) Shutdown() error {
	log.Info("server start to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return this.Server.Shutdown(ctx)
}

func (this *Server) EnableMetrics(path string) *Server {
	this.GET(path, Metrics)
	return this
}

func (this *Server) AddRouter(router *Router) *Server {
	this.recursiveAddRouter(&this.Engine.RouterGroup, router)
	return this
}

func (this *Server) recursiveAddRouter(group *gin.RouterGroup, router *Router) {
	if len(router.Children) > 0 {
		child := group.Group(router.Path)
		for _, r := range router.Children {
			this.recursiveAddRouter(child, r)
		}
	}
	if router.Handler != nil {
		this.addRouter(group, router)
	}
}

func (this *Server) addRouter(group *gin.RouterGroup, router *Router) {
	var addRouter func(path string, handler ... gin.HandlerFunc) gin.IRoutes
	switch router.Method {
	case http.MethodGet:
		addRouter = group.GET
	case http.MethodPost:
		addRouter = group.POST
	case http.MethodPut:
		addRouter = group.PUT
	case http.MethodDelete:
		addRouter = group.DELETE
	case http.MethodPatch:
		addRouter = group.PATCH
	case http.MethodHead:
		addRouter = group.HEAD
	case http.MethodOptions:
		addRouter = group.OPTIONS
	case HttpMethodAny:
		addRouter = group.Any
	default:
		panic(errors.New("unknown method " + router.Method))
	}
	addRouter(router.Path, router.Handler)
}
