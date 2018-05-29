package dubbo

import (
	"bytes"
	"encoding/binary"

	"dubbo-mesh/json"
)

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
	data = bytes.Join([][]byte{dubbo, path, version, method, paramType, inv.Args, attach}, ParamSeparator)
	return data
}
