//一些基本的辅助类
package util

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WaitForExitSign() {
	c := make(chan os.Signal, 1)
	//结束，收到ctrl+c 信号
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-c
}

//休眠n秒
func Sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}
