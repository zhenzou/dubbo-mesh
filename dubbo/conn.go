package dubbo

import (
	"net"

	"dubbo-mesh/log"
)

func NewPool(addr string) *Pool {
	log.Debug("addr:", addr)
	return &Pool{addr}
}

type Pool struct {
	addr string
}

func (this *Pool) Get() net.Conn {
	conn, err := net.Dial("tcp", this.addr)
	if err != nil {
		panic(err)
	}
	return conn
}

func (this *Pool) Put(conn net.Conn) {
	conn.Close()
}
