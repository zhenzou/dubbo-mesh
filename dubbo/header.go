package dubbo

import (
	"sync"
)

const (
	HeaderLength      = 16
	FlagRequest  byte = 0x80
	FlagTwoWay   byte = 0x40
	FlagEvent    byte = 0x20

	StatusOk                             = 20
	StatusClientTimeOut                  = 30
	StatusServerTimeout                  = 31
	StatusBadReq                         = 40
	StatusBadResp                        = 50
	StatusServiceNotFound                = 60
	StatusServiceError                   = 70
	StatusServerError                    = 80
	StatusClientError                    = 90
	StatusServerThreadPoolExhaustedError = 100
)

var (
	Magic      = []byte{0xda, 0xbb}
	headerPool sync.Pool
)

func init() {
	headerPool = sync.Pool{New: NewHeader}
}

func NewHeader() interface{} {
	header := make([]byte, HeaderLength)
	header[0] = Magic[0]
	header[1] = Magic[1]
	return Header(header)
}

type Header []byte

// 判断是否请求，如果是请求则忽略
func (this Header) Req() bool {
	return this[2]&FlagRequest != 0
}

// 通过Header获取payload的长度
func (this Header) DataLen() int {
	return Bytes2Int(this[12:])
}

// 通过Header获取payload的长度
func (this Header) RequestId() int64 {
	return Bytes2Int64(this[4:12])
}

// 返回值状态
func (this Header) Status() int {
	return int(this[3])
}

func (this Header) Bytes() []byte {
	return []byte(this)
}
