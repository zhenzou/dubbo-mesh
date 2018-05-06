package dubbo

import (
	"dubbo-mesh/web"
	"github.com/gin-gonic/gin"
	"dubbo-mesh/log"
	"net/http"
)

type Config struct {
	DubboAddr  string
	ServerPort int
	Etcd       string
}

var (
	cfg    *Config
	client *Client
)

func invoke(ctx *gin.Context) {
	interfaceName := ctx.PostForm("interface")
	method := ctx.PostForm("method")
	paramType := ctx.PostForm("parameterTypesString")
	param := ctx.PostForm("parameter")
	log.Debugf("%s/%s/%s/%s", interfaceName, method, paramType, param)
	resp, err := client.Invoke(interfaceName, method, paramType, param)
	if err != nil {
		ctx.AbortWithStatus(500)
	} else {
		ctx.Data(http.StatusOK, web.MIMEApplicationJSONCharsetUTF8, resp)
	}
}

func NewServer(cfg *Config) *web.Server {
	server := web.NewServer()
	client = NewClient(cfg.DubboAddr)

	server.POST("", invoke)

	return server
}
