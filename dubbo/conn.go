package dubbo

import (
	"net"
	"errors"

	"dubbo-mesh/log"
)

var (
	ReadBeforeRequestError = errors.New("")
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

func NewPool(dubboAddr string) *Pool {
	log.Debug("dubboAddr:", dubboAddr)
	pool := &Pool{addr: dubboAddr, ch: make(chan *Conn, 256)}
	return pool
}

type Pool struct {
	addr string
	ch   chan *Conn
}

func (this *Pool) new() *Conn {
	conn, err := net.Dial("tcp", this.addr)
	if err != nil {
		panic(err)
	}
	log.Info("new:", conn.LocalAddr())

	return &Conn{Conn: conn}
}
func (this *Pool) Get() *Conn {
	select {
	case conn := <-this.ch:
		return conn
	default:
		return this.new()
	}
}

// TODO POOl
func (this *Pool) Put(conn *Conn) {
	select {
	case this.ch <- conn:
	default:
		log.Info("close:", conn.LocalAddr())
		conn.Close()
	}
}

func (this *Pool) Shutdown() {
	for conn := range this.ch {
		conn.Close()
	}
}
