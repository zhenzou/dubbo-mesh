package main

import (
	"dubbo-mesh/dubbo"
	"dubbo-mesh/util"
	"dubbo-mesh/derror"
	"net/http"
)

func main() {
	cfg := &dubbo.Config{
		DubboAddr:  "127.0.0.1:20880",
		ServerPort: 20000,
		Etcd:       "http://127.0.0.1:2379",
	}
	server := dubbo.NewServer(cfg)

	go func() {
		if err := server.Run(cfg.ServerPort); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	util.WaitForExitSign()

	derror.Warn(server.Shutdown())
}
