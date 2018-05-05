package web

import (
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"

	"dubbo-mesh/log"
	"dubbo-mesh/util"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				request, _ := httputil.DumpRequest(ctx.Request, false)
				log.Error(util.BytesToString(request), err)
				ctx.AbortWithStatus(500)
			}
		}()
		ctx.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path

		// Process request
		ctx.Next()

		// Log only when path is not being skipped
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()

		comment := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Infof("\033[36;1m[GIN]\033[0m  \033[36;1m[%.2fms]\033[0m | %3d | %15s |%s  %s  %s", float64(latency.Nanoseconds()/time.Millisecond.Nanoseconds()), statusCode, clientIP, method, path, comment)

	}
}
