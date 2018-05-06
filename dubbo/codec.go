package dubbo

import (
	"bytes"

	"dubbo-mesh/json"
	"dubbo-mesh/log"
	"encoding/binary"
)

const (
	HeaderLength      = 16
	FlagRequest  byte = 0x80
	FlagTwoWay   byte = 0x40
	FlagEvent    byte = 0x20
)

var (
	Magic = []byte{0xda, 0xbb}
)

type Codec struct {
}

func EncodeInt16(encode []byte, i int16, offset ...int) []byte {
	off := 0
	if len(offset) > 0 {
		off = offset[0]
	}
	encode[off+1] = byte(i)
	encode[off+0] = byte(i >> 8)
	return encode
}

func EncodeInt64(encode []byte, i int64, offset ...int) []byte {
	off := 0
	if len(offset) > 0 {
		off = offset[0]
	}
	encode[off+7] = byte(i)
	encode[off+6] = byte(i >> 8)
	encode[off+5] = byte(i >> 16)
	encode[off+4] = byte(i >> 24)
	encode[off+3] = byte(i >> 32)
	encode[off+2] = byte(i >> 40)
	encode[off+1] = byte(i >> 48)
	encode[off+0] = byte(i >> 56)
	return encode
}

func EncodeInt(encode []byte, i int, offset ...int) []byte {
	off := 0
	if len(offset) > 0 {
		off = offset[0]
	}
	encode[off+3] = byte(i)
	encode[off+2] = byte(i >> 8)
	encode[off+1] = byte(i >> 16)
	encode[off+0] = byte(i >> 24)
	return encode
}

func Bytes2Int(bytes []byte) int {
	return int(binary.BigEndian.Uint32(bytes))
}

func Bytes2Int64(bytes []byte) int64 {
	return int64(binary.BigEndian.Uint64(bytes))
}

func EncodeInvocation(inv *Invocation) []byte {
	var data []byte
	dubbo, _ := json.Marshal(inv.Attach["dubbo"])
	path, _ := json.Marshal(inv.Attach["path"])
	version, _ := json.Marshal(inv.Attach["version"])
	method, _ := json.Marshal(inv.Method)
	paramType, _ := json.Marshal(inv.ParamType)
	attach, _ := json.Marshal(inv.Attach)
	data = bytes.Join([][]byte{dubbo, path, version, method, paramType, inv.Args, attach}, []byte("\n"))
	//data = append(data, dubbo...)
	//data = append(data, path...)
	//data = append(data, version...)
	//data = append(data, method...)
	//data = append(data, paramType...)
	log.Debug("data:", string(data))
	return data
}

func Encode(req *Request) []byte {
	header := make([]byte, HeaderLength)
	header[0] = Magic[0]
	header[1] = Magic[1]
	header[2 ] = FlagRequest | 6
	if req.TwoWay {
		header[2] |= FlagTwoWay
	}
	if req.Event {
		header[2] |= FlagEvent
	}
	EncodeInt64(header, req.Id, 4)
	data := EncodeInvocation(req.Data.(*Invocation))
	EncodeInt(header, len(data), 12)
	buf := bytes.NewBuffer(make([]byte, 0, HeaderLength+len(data)))
	buf.Write(header)
	buf.Write(data)
	return buf.Bytes()
}
