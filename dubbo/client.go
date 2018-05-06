package dubbo

import (
	"dubbo-mesh/json"
)

type Process func(request *Request) (resp *Response, err error)

func NewClient(addr string) *Client {
	client := &Client{pool: NewPool(addr)}
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

	return this.process(req)
}

func (this *Client) getConn() *Conn {
	return this.pool.Get()
}

func (this *Client) closeConn(conn *Conn) {
	this.pool.Put(conn)
}

// TODO Retry
func (this *Client) defaultProcess(request *Request) (resp *Response, err error) {
	conn := this.getConn()
	defer this.closeConn(conn)
	err = conn.WriteRequest(request)
	if err != nil {
		return
	}
	return conn.ReadResponse()
}
