package mesh

import (
	"net"
	"fmt"
	"sync"
	"strings"
	"bytes"

	"github.com/xtaci/kcp-go"

	"dubbo-mesh/log"
	"dubbo-mesh/registry"
	"dubbo-mesh/dubbo"
	"dubbo-mesh/util"
)

func NewKcpClient() Client {
	return &KcpClient{
		pool: make(map[*registry.Endpoint]*Pool),
	}
}

// TODO 与TCP整合
type KcpClient struct {
	sync.Mutex
	pool map[*registry.Endpoint]*Pool
}

func (this *KcpClient) newPool(addr string) *Pool {
	return NewPool(200, func() (net.Conn, error) {
		conn, err := kcp.Dial(addr)
		if err != nil {
			return nil, err
		}
		log.Info("new ", conn.LocalAddr())

		return conn, nil
	})
}

func (this *KcpClient) Invoke(endpoint *registry.Endpoint, inv *Invocation) ([]byte, error) {
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
	defer pool.Put(conn)
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

func NewKcpServer(port int, dubbo *dubbo.Client) Server {
	return &KcpServer{client: dubbo, addr: fmt.Sprintf(":%d", port)}
}

type KcpServer struct {
	addr     string
	listener net.Listener
	client   *dubbo.Client
}

func (this *KcpServer) Invocations() <-chan Invocation {
	panic("implement me")
}

func (this *KcpServer) Run() error {
	listener, err := kcp.Listen(this.addr)
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

func (this *KcpServer) handle(conn net.Conn) error {
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

func (this *KcpServer) Shutdown() error {
	return this.listener.Close()
}
