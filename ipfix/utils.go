package ipfix

import (
	"bytes"
	"encoding/binary"
)

//Int16ToBytes 将int16转换为2个byte
func Int16ToBytes(data int16) (byte, byte) {
	return byte(data >> 8), byte(data & ((1 << 8) - 1))
}

//Int32ToBytes 将int32转换为4个byte
func Int32ToBytes(data int32) (byte, byte, byte, byte) {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	bt := bytebuf.Bytes()
	return bt[0],bt[1],bt[2],bt[3]
}