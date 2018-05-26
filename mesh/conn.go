package mesh

import (
	"net"
	"time"
	"context"
	"errors"
	"sync/atomic"

	"dubbo-mesh/derror"
	"dubbo-mesh/log"
)

var (
	PoolShutdownError = errors.New("pool shutdown")
	ErrorResp         = []byte("error")
)

func NewPool(max int, new func() (net.Conn, error)) *Pool {
	return &Pool{New: new, ch: make(chan net.Conn, max)}
}

type Pool struct {
	addr  string
	ch    chan net.Conn
	count uint32
	New   func() (net.Conn, error)
}

func (this *Pool) Get() (net.Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		atomic.AddUint32(&this.count, 1)
		conn, err := this.New()
		if err == nil {
			log.Infof("new %d %s", this.count, conn.RemoteAddr())
		}
		return conn, err
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
