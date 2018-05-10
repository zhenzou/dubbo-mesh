package dubbo

import (
	"net"
	"errors"
	"context"
	"time"

	"dubbo-mesh/log"
	"dubbo-mesh/derror"
)

var (
	ReadBeforeRequestError = errors.New("")
	PoolShutdownError      = errors.New("pool shutdown")
	ParamSeparator         = []byte("\n")
)

type Conn struct {
	net.Conn
	send bool
}

func (this *Conn) WriteRequest(req *Request) (err error) {
	header := headerPool.Get().(Header).Bytes()
	if req.TwoWay {
		header[2] |= FlagTwoWay
	}
	if req.Event {
		header[2] |= FlagEvent
	}
	EncodeInt64(header, req.Id, 4)
	data := EncodeInvocation(req.Data.(*Invocation))
	EncodeInt(header, len(data), 12)

	payload := make([]byte, HeaderLength+len(data))

	copy(payload, header)
	copy(payload[HeaderLength:], data)
	_, err = this.Write(payload)
	if err != nil {
		return
	}
	this.send = true
	return
}

func (this *Conn) Close() (err error) {
	log.Info("close ", this.LocalAddr())
	return this.Conn.Close()
}

func (this *Conn) ReadResponse() (resp *Response, err error) {
	if !this.send {
		err = ReadBeforeRequestError
		return
	}
	header := headerPool.Get().(Header)
	var _ int
	_, err = this.Read(header)
	if err != nil {
		return
	}
	length := header.DataLen()
	data := make([]byte, length)
	_, err = this.Read(data)
	if err != nil {
		return
	}
	resp = NewResponse(header.Status(), header.RequestId(), data)
	return
}

func NewPool(max int, dubboAddr string) *Pool {
	log.Debug("dubboAddr ", dubboAddr)
	pool := &Pool{addr: dubboAddr, ch: make(chan *Conn, max)}
	return pool
}

type Pool struct {
	addr string
	ch   chan *Conn
}

func (this *Pool) new() (*Conn, error) {
	conn, err := net.Dial("tcp", this.addr)
	if err != nil {
		return nil, err
	}
	log.Info("new ", conn.LocalAddr())

	return &Conn{Conn: conn}, nil
}
func (this *Pool) Get() (*Conn, error) {
	select {
	case conn, more := <-this.ch:
		if !more {
			return nil, PoolShutdownError
		}
		return conn, nil
	default:
		return this.new()
	}
}

// TODO POOl
func (this *Pool) Put(conn *Conn) {
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
