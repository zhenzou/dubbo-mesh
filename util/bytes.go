package util

import (
	"encoding/binary"
)

func Int2Bytes(i int) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(i))
	return data
}

func Bytes2Int(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}
