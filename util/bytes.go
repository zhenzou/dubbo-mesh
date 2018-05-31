package util

import (
	"encoding/binary"
)

func Int2Bytes(v int) []byte {
	bs := make([]byte, 4)
	bs[0] = byte(v >> 24)
	bs[1] = byte(v >> 16)
	bs[2] = byte(v >> 8)
	bs[3] = byte(v)
	return bs
}

func Int64ToBytes(i int64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(i))
	return data
}

func Bytes2Int(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

func Bytes2Int64(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}
