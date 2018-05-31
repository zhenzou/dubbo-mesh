package dubbo

import (
	"sync"

	"dubbo-mesh/util"
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
	Magic   = []byte{0xda, 0xbb}
	headers = sync.Pool{New: NewHeader}
)

func NewHeader() interface{} {
	header := make([]byte, HeaderLength)
	header[0] = Magic[0]
	header[1] = Magic[1]
	return Header(header)
}

type Header []byte

// 判断是否请求
func (this Header) Req() bool {
	return this[2]&FlagRequest != 0
}

// 设置header是请求
func (this Header) SetReq() {
	this[2] = FlagRequest | 6
}

func (this Header) SetTwoWay() {
	this[2] |= FlagTwoWay
}

func (this Header) SetEvent() {
	this[2] |= FlagEvent
}

// 通过Header获取payload的长度
func (this Header) Len() int {
	return util.Bytes2Int(this[12:])
}

// 设置请求payload的长度
func (this Header) SetLen(length int) {
	this[15] = byte(length)
	this[14] = byte(length >> 8)
	this[13] = byte(length >> 16)
	this[12] = byte(length >> 24)
}

// 获取响应的ReqID
func (this Header) ReqId() int64 {
	return util.Bytes2Int64(this[4:12])
}

// 设置请求头的ReqID
func (this Header) SetReqId(id int64) {
	this[11] = byte(id)
	this[10] = byte(id >> 8)
	this[9] = byte(id >> 16)
	this[8] = byte(id >> 24)
	this[7] = byte(id >> 32)
	this[6] = byte(id >> 40)
	this[5] = byte(id >> 48)
	this[4] = byte(id >> 56)
}

// 返回值状态
func (this Header) Status() int {
	return int(this[3])
}

func (this Header) Bytes() []byte {
	return []byte(this)
}
