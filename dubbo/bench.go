// +build prod

package dubbo

import (
	"sync/atomic"
	"dubbo-mesh/log"
)

// 为了bench优化，在mesh端已经限制连接数，这里不会超过了
func (this *Pool) Get() (*Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		conn, err := this.new()
		if err == nil {
			count := atomic.AddUint32(&this.count, 1)
			log.Infof("dubbo %d %s", count, conn.RemoteAddr())
		}
		return conn, err
	}
}
