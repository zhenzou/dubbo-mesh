package dubbo

import (
	"sync"
)

const (
	HeaderLength      = 16
	FlagRequest  byte = 0x80
	FlagTwoWay   byte = 0x40
	FlagEvent    byte = 0x20
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
	header[2 ] = FlagRequest | 6
	return Header(header)
}

type Header []byte

// 通过Header获取payload的长度
func (this Header) DataLen() int {
	return Bytes2Int(this[12:])
}

// 通过Header获取payload的长度
func (this Header) RequestId() int64 {
	return Bytes2Int64(this[4:12])
}

func (this Header) Bytes() []byte {
	return []byte(this)
}
