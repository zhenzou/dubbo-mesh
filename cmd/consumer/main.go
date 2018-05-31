package main

import (
	"fmt"
	"net/http"

	"dubbo-mesh/cmd"
	"dubbo-mesh/derror"
	"dubbo-mesh/log"
	"dubbo-mesh/sidecar"
	"dubbo-mesh/util"
)

func main() {

	var (
		cfg    *sidecar.Config
		server *sidecar.Consumer
	)

	cmd.Run(
		func() error {
			cfg = &sidecar.Config{
				DubboAddr:  fmt.Sprintf("127.0.0.1:%d", cmd.DubboPort),
				ServerPort: cmd.Port,
				Etcd:       cmd.Etcd,
				Service:    cmd.Service,
				Balancer:   sidecar.LB_WLActive,
				Server:     1,
			}
			log.Info(util.ToJsonStr(cfg))
			return nil
		},
		func() error {
			server = sidecar.NewConsumer(cfg)
			if err := server.Run(); err != http.ErrServerClosed {
				log.Panic(err)
			}
			return nil
		},
		func() error {
			if server != nil {
				derror.Panic(server.Shutdown())
			}
			return nil
		},
	)
}
