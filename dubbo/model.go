package dubbo

import (
	"sync/atomic"
	"errors"
	"bytes"

	"dubbo-mesh/util"
)

var (
	nextId int64
)

const (
	ResponseNullValue     = 2
	ResponseValue         = 1
	ResponseWithException = 0
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

func NewResponse(status int, id int64, data []byte) *Response {
	return &Response{status, id, data, ""}
}

type Response struct {
	Status   int
	ReqId    int64  // requestId
	Payload  []byte // 数据
	ErrorMsg string
}

func (this *Response) Error() error {
	if this.Status == StatusOk {
		return nil
	}
	return errors.New(string(this.Body()))
}

// 返回body，错误消息，或者返回值
// 为了比赛的 case 优化
func (this *Response) Body() []byte {
	split := bytes.Split(this.Payload, ParamSeparator)

	if this.Status == StatusOk {
		//data = bytes.Join(split[1:len(split)-1], ParamSeparator)
		return util.TrimCR(split[1])
	}
	return split[0]
}
