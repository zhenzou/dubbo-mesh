package mesh

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"dubbo-mesh/dubbo"
	"dubbo-mesh/json"
	"dubbo-mesh/log"
	"dubbo-mesh/registry"
	"dubbo-mesh/util"
)

// 为了优化，一般来说应该要动态的初始化连接池
func NewTcpClient(endpoints []*registry.Endpoint) Client {
	client := &TcpClient{
		pool: make(map[*registry.Endpoint]*Pool),
	}
	for _, endpoint := range endpoints {
		client.pool[endpoint] = client.newPool(endpoint.Addr())
	}
	return client
}

type TcpClient struct {
	sync.Mutex
	pool map[*registry.Endpoint]*Pool
}

func (this *TcpClient) newPool(addr string) *Pool {
	return NewPool(200, func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	})
}

func (this *TcpClient) Invoke(endpoint *registry.Endpoint, inv *Invocation) ([]byte, error) {

	var (
		pool *Pool
		ok   bool
	)
	// DCL
	if pool, ok = this.pool[endpoint]; !ok {
		return nil, errors.New("not registered")
	}
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}
	data := inv.Data()
	conn.Write(util.Int2Bytes(len(data)))
	conn.Write(data)
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
	lengthBuf := make([]byte, 4)
	for {
		_, err := conn.Read(lengthBuf)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		length := util.Bytes2Int(lengthBuf)
		buf := buf[:length]
		_, err = conn.Read(buf)
		if err != nil {
			log.Warn(err.Error())
			conn.Write(ErrorResp)
			break
		}

		inv := NewInv()
		// 忽略错误处理
		json.Unmarshal(buf, inv)

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
