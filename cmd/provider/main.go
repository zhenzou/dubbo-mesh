package main

import (
	"fmt"

	"dubbo-mesh/cmd"
	"dubbo-mesh/log"
	"dubbo-mesh/sidecar"
	"dubbo-mesh/util"
)

func main() {

	var (
		cfg    *sidecar.Config
		server *sidecar.Provider
	)

	cmd.Run(
		func() error {
			cfg = &sidecar.Config{
				DubboAddr:  fmt.Sprintf("127.0.0.1:%d", cmd.DubboPort),
				ServerPort: cmd.Port,
				Etcd:       cmd.Etcd,
				Service:    cmd.Service,
			}
                        log.Debug("cfg:",util.ToJsonStr(cfg))
			return nil
		},
		func() error {
			server = sidecar.NewProvider(cfg)
			if err := server.Run(); err != nil {
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
