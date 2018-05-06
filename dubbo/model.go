package dubbo

import "sync/atomic"

var (
	nextId int64
)

type Invocation struct {
	ParamType string
	Method    string
	Args      []byte
	Attach    map[string]interface{}
}

type Request struct {
	DubboVersion string
	Version      string
	Id           int64
	Interface    string
	Method       string
	ParamType    string // 参数类型
	TwoWay       bool
	Event        bool
	Data         interface{}
}

func NewRequest(version, interfaceName, method, paramType string, data interface{}) *Request {
	id := atomic.AddInt64(&nextId, 1)
	return &Request{
		DubboVersion: "2.6.0",
		Version:      version,
		Id:           id,
		Interface:    interfaceName,
		Method:       method,
		TwoWay:       true,
		ParamType:    paramType,
		Data:         data,
	}
}

func NewResponse(id int64, data []byte) *Response {
	return &Response{id, data}
}

type Response struct {
	ReqId   int64  // requestId
	Payload []byte // 数据
}
