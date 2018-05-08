package main

import (
	"fmt"
	"net/http"

	"dubbo-mesh/app"
	"dubbo-mesh/log"
	"dubbo-mesh/sidecar"
)

func main() {

	var (
		cfg    *sidecar.Config
		server *sidecar.Provider
	)

	app.Run(
		func() error {
			cfg = &sidecar.Config{
				DubboAddr:  fmt.Sprintf("127.0.0.1:%d", app.DubboPort),
				ServerPort: app.Port,
				Etcd:       app.Etcd,
				Service:    app.Service,
			}
			return nil
		},
		func() error {
			server = sidecar.NewProvider(cfg)
			if err := server.Run(); err != http.ErrServerClosed {
				log.Panic(err)
			}
			return nil
		},
		func() error {
			server.Shutdown()

			return nil
		},
	)
}
