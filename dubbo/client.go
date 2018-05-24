package dubbo

import (
	"dubbo-mesh/json"
)

type Process func(conn *Conn, request *Request) (resp *Response, err error)

func NewClient(addr string, size int) *Client {
	client := &Client{pool: NewPool(size, addr)}
	client.init()
	return client
}

type Client struct {
	pool    *Pool
	process Process
}

func (this *Client) init() {
	this.process = this.defaultProcess
}

func (this *Client) Invoke(interfaceName, method, paramType, param string) (resp *Response, err error) {
	data, _ := json.Marshal(param)

	invocation := Invocation{
		Attach:    map[string]interface{}{"path": interfaceName, "dubbo": "2.0.1"},
		Method:    method,
		ParamType: paramType,
		Args:      data,
	}

	req := NewRequest("2.0.0", interfaceName, method, paramType, &invocation)

	conn := this.getConn()
	resp, err = this.process(conn, req)
	if err != nil {
		conn.Close()
		return
	}
	this.closeConn(conn)
	return
}

func (this *Client) Shutdown() {
	this.pool.Shutdown()
}

func (this *Client) getConn() *Conn {
	// 暂时忽略错误
	conn, _ := this.pool.Get()
	return conn
}

func (this *Client) closeConn(conn *Conn) {
	this.pool.Put(conn)
}

// TODO Retry
func (this *Client) defaultProcess(conn *Conn, request *Request) (resp *Response, err error) {
	err = conn.WriteRequest(request)
	if err != nil {
		return
	}
	return conn.ReadResponse()
}
