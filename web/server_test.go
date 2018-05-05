package web

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func handler(ctx *gin.Context) {

}

func TestEngine_AddRouter(t *testing.T) {
	engine := NewServer()
	type args struct {
		router *Router
	}

	router := Router{
		Method: http.MethodGet,
		Path:   "L1",
		Children: []*Router{
			{http.MethodGet, "L21", handler, []*Router{
				{http.MethodGet, "L31", handler, nil},
				{http.MethodGet, "L32", handler, nil},
			},
			},
			{http.MethodGet, "L22", nil, []*Router{
				{http.MethodGet, "L33", handler, nil},
				{http.MethodGet, "L34", handler, nil},
			},
			},
		},
	}
	engine.AddRouter(&router)
	engine.Routes()
}
