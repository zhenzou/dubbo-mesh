package mesh

import (
	"net"
	"fmt"
	"sync"
	"strings"
	"bytes"
	"io"

	"dubbo-mesh/log"
	"dubbo-mesh/registry"
	"dubbo-mesh/dubbo"
	"dubbo-mesh/util"
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

func (this *TcpClient) newPool(addr string) *Pool {
	return NewPool(256, func() (net.Conn, error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		log.Info("new ", conn.LocalAddr())

		return conn, nil
	})
}

func (this *TcpClient) Invoke(endpoint *registry.Endpoint, inv *Invocation) ([]byte, error) {
	var (
		pool *Pool
		ok   bool
	)
	// DCL
	if pool, ok = this.pool[endpoint]; !ok {
		this.Lock()
		if pool, ok = this.pool[endpoint]; !ok {
			pool = this.newPool(endpoint.String())
			this.pool[endpoint] = pool
		}
		this.Unlock()
	}
	conn, _ := pool.Get()
	data := strings.Join([]string{inv.Interface, inv.Method, inv.ParamType, inv.Param}, "\n")

	conn.Write(util.Int2Bytes(len(data)))
	conn.Write(util.StringToBytes(data))
	//
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		log.Warn(err.Error())
		conn.Close()
		return nil, err
	}
	pool.Put(conn)
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
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Warn("err:", err.Error())
			continue
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
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		length := util.Bytes2Int(l)
		buf := buf[:length]
		_, err = conn.Read(buf)
		if err != nil {
			log.Warn(err.Error())
			break
		}
		split := bytes.Split(buf, []byte("\n"))

		resp, err := this.client.Invoke(util.BytesToString(split[0]), util.BytesToString(split[1]), util.BytesToString(split[2]), util.BytesToString(split[3]))
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
