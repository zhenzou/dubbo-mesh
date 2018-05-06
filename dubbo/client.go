package dubbo

import (
	"bytes"

	"dubbo-mesh/json"
	"dubbo-mesh/log"
)

func NewClient(addr string) *Client {
	return &Client{NewPool(addr)}
}

type Client struct {
	pool *Pool
}

func (this *Client) Invoke(interfaceName, method, paramType, param string) (resp []byte, err error) {
	data, _ := json.Marshal(param)

	invocation := Invocation{
		Attach:    map[string]interface{}{"path": interfaceName, "dubbo": "2.0.1"},
		Method:    method,
		ParamType: paramType,
		Args:      data,
	}

	req := NewRequest("2.0.0", interfaceName, method, paramType, &invocation)
	conn := this.pool.Get()
	defer this.pool.Put(conn)
	data = Encode(req)
	_, err = conn.Write(data)
	if err != nil {
		log.Warn(err.Error())
	}
	header := make([]byte, HeaderLength)
	n, err := conn.Read(header)
	if err != nil {
		panic(err)
	}
	log.Debug("read:", n)
	id := Bytes2Int64(header[4:12])
	log.Debug("id:", id)
	dataLen := Bytes2Int(header[12:])
	log.Debug("len:", dataLen)
	data = make([]byte, dataLen)
	n, err = conn.Read(data)
	if err != nil {
		panic(err)
	}
	log.Debug("read:", n)
	log.Debug(string(data))
	split := bytes.Split(data, []byte("\n"))
	resp = split[1]
	return
}
