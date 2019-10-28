package ipfix

import (
	"bytes"
	"encoding/binary"
)

//IntToBytes 将int转换为2个byte
func IntToBytes(data int) (byte, byte) {
	return byte(data >> 8), byte(data & ((1 << 8) - 1))
}

//Int64ToBytes 将int64转换为4个byte
func Int64ToBytes(data int64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}