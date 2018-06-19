package mesh

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"dubbo-mesh/derror"
	"dubbo-mesh/log"
)

var (
	PoolShutdownError = errors.New("pool shutdown")
	ErrorResp         = []byte("error")
)

func NewPool(max int, new func() (net.Conn, error)) *Pool {
	return &Pool{New: new, max: max, ch: make(chan net.Conn, max)}
}

type Pool struct {
	addr  string
	ch    chan net.Conn
	max   int
	count int
	mtx   sync.Mutex
	New   func() (net.Conn, error)
}

// 一定要使用Put放回来
func (this *Pool) Get() (net.Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		this.mtx.Lock()
		if this.count < this.max {
			conn, err := this.New()
			if err == nil {
				this.count += 1
				log.Infof("mesh %d %s", this.count, conn.RemoteAddr())
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

func (this *Pool) Put(conn net.Conn) {
	select {
	case this.ch <- conn:
	default:
		log.Info("close:", conn.RemoteAddr())
		conn.Close()
	}
}

func (this *Pool) Shutdown() error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	close(this.ch)
	done := make(chan struct{})
	go func() {
		for conn := range this.ch {
			derror.Warn(conn.Close())
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
