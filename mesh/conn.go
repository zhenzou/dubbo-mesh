package mesh

import (
	"net"
	"time"
	"context"
	"errors"

	"dubbo-mesh/derror"
)

var (
	PoolShutdownError = errors.New("pool shutdown")
	ErrorResp         = []byte("error")
)

func NewPool(max int, new func() (net.Conn, error)) *Pool {
	return &Pool{New: new, ch: make(chan net.Conn, max)}
}

type Pool struct {
	addr string
	ch   chan net.Conn
	New  func() (net.Conn, error)
}

func (this *Pool) Get() (net.Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		return this.New()
	}
}

// TODO POOl
func (this *Pool) Put(conn net.Conn) {
	select {
	case this.ch <- conn:
	default:
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
