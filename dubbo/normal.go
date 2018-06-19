// +build !prod

package dubbo

import (
	"dubbo-mesh/log"
)

// 正常情况下的连接池Get方法
func (this *Pool) Get() (*Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		this.mtx.Lock()
		if this.count < this.max {
			conn, err := this.new()
			if err == nil {
				this.count += 1
				log.Infof("dubbo %d %s", this.count, conn.RemoteAddr())
			}
			this.mtx.Unlock()
			return conn, err
		} else {
			this.mtx.Unlock()
			select {
			case conn, more := <-this.ch:
				if !more {
					return nil, PoolShutdownError
				}
				return conn, nil
			}
		}
	}
}
