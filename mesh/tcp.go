package mesh

import (
	"time"
	"net"
	"errors"
	"context"
	"fmt"
	"sync"
	"strings"
	"bytes"

	"dubbo-mesh/derror"
	"dubbo-mesh/log"
	"dubbo-mesh/registry"
	"dubbo-mesh/dubbo"
	"dubbo-mesh/util"
)

var (
	PoolShutdownError = errors.New("pool shutdown")
	ErrorResp         = []byte("error")
)

func NewTcpClient() Client {
	return &TcpClient{
		pool: make(map[*registry.Endpoint]*Pool),
	}
}

type TcpClient struct {
	sync.Mutex
	pool map[*registry.Endpoint]*Pool
}

func (this *TcpClient) Invoke(endpoint *registry.Endpoint, inv *Invocation) ([]byte, error) {
	var conn net.Conn
	// DCL
	if pool, ok := this.pool[endpoint]; !ok {
		this.Lock()
		if pool, ok = this.pool[endpoint]; !ok {
			pool = NewPool(200, endpoint.String())
			this.pool[endpoint] = pool
		}
		this.Unlock()
		conn, _ = pool.Get()
		defer pool.Put(conn)
	} else {
		conn, _ = pool.Get()
		defer pool.Put(conn)
	}
	data := strings.Join([]string{inv.Interface, inv.Method, inv.ParamType, inv.Param}, "\n")

	conn.Write(util.Int2Bytes(len(data)))
	conn.Write(util.StringToBytes(data))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Warn(err.Error())
		conn.Close()
		return nil, err
	}
	return buf[:n], nil
}

func NewTcpServer(port int, dubbo *dubbo.Client) Server {
	return &TcpServer{client: dubbo, addr: fmt.Sprintf(":%d", port)}
}

type TcpServer struct {
	addr     string
	listener net.Listener
	client   *dubbo.Client
}

func (this *TcpServer) Invocations() <-chan Invocation {
	panic("implement me")
}

func (this *TcpServer) Run() error {
	listener, err := net.Listen("tcp", this.addr)
	if err != nil {
		return err
	}
	this.listener = listener
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Warn(err.Error())
			break
		}
		go this.handle(conn)
	}
	return err
}

func (this *TcpServer) handle(conn net.Conn) error {
	buf := make([]byte, 2048)
	l := make([]byte, 4)
	for {
		_, err := conn.Read(l)
		if err != nil {
			log.Warn(err.Error())
			break
		}
		length := util.Bytes2Int(l)
		log.Debug("length:", length)

		buf := buf[:length]
		_, err = conn.Read(buf)
		if err != nil {
			log.Warn(err.Error())
			break
		}
		inv := &Invocation{}
		split := bytes.Split(buf, []byte("\n"))
		inv.Interface = util.BytesToString(split[0])
		inv.Method = util.BytesToString(split[1])
		inv.ParamType = util.BytesToString(split[2])
		inv.Param = util.BytesToString(split[3])

		resp, err := this.client.Invoke(inv.Interface, inv.Method, inv.ParamType, inv.Param)
		if err != nil {
			log.Warn(err.Error())
			conn.Write(ErrorResp)
		} else if resp.Error() != nil {
			log.Warn(resp.Error().Error())
			conn.Write(ErrorResp)
		} else {
			conn.Write(resp.Body())
		}
	}
	return conn.Close()
}

func (this *TcpServer) Shutdown() error {
	return this.listener.Close()
}

func NewPool(max int, addr string) *Pool {
	log.Debug("addr ", addr)
	pool := &Pool{addr: addr, ch: make(chan net.Conn, max)}
	return pool
}

type Pool struct {
	addr string
	ch   chan net.Conn
}

func (this *Pool) new() (net.Conn, error) {
	conn, err := net.Dial("tcp", this.addr)
	if err != nil {
		return nil, err
	}
	log.Info("new ", conn.LocalAddr())

	return conn, nil
}
func (this *Pool) Get() (net.Conn, error) {
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
