package ipfix

import (
	"bytes"
	"encoding/binary"
)

//Int16ToBytes 将int16转换为2个byte
func Int16ToBytes(data int16) (byte, byte) {
	return byte(data >> 8), byte(data & ((1 << 8) - 1))
}

//IntToBytes 将int转换为2个byte
func IntToBytes(data int) (byte, byte) {
	return byte(data >> 8), byte(data & ((1 << 8) - 1))
}

//Int32ToBytes 将int32转换为4个byte
func Int32ToBytes(data int32) (byte, byte, byte, byte) {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	bt := bytebuf.Bytes()
	return bt[0],bt[1],bt[2],bt[3]
}

//Int64ToBytes 将int64转换为4个byte,只取最后4位
func Int64ToBytes(data int64) (byte, byte, byte, byte) {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	bt := bytebuf.Bytes()
	return bt[4],bt[5],bt[6],bt[7]
}

//checkSum 计算ip头和udp头的校验算法
func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var ans = uint16(^sum)
	return ans
}
