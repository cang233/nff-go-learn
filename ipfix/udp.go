package ipfix

import "net"

//GenerateUDPHeader 生成udp报文头
func GenerateUDPHeader(srcIP, dstIP net.IP, srcPort, dstPort int,data []byte) *[]byte {
	//填充udp首部
	//udp伪首部
	udph := make([]byte, 20)
	//源ip地址
	udph[0], udph[1], udph[2], udph[3] = srcIP[12], srcIP[13], srcIP[14], srcIP[15]
	//目的ip地址
	udph[4], udph[5], udph[6], udph[7] = dstIP[12], dstIP[13], dstIP[14], dstIP[15]
	//协议类型
	udph[8], udph[9] = 0x00, 0x11
	//udp头长度
	udph[10], udph[11] = IntToBytes(len(data)+8)
	//下面开始就真正的udp头部
	//源端口号
	udph[12], udph[13] = IntToBytes(srcPort)
	//目的端口号，改成2055后wireshark就自动解析成cflow协议了。
	udph[14], udph[15] = IntToBytes(dstPort)
	//udp头长度
	udph[16], udph[17] = IntToBytes(len(data)+8)
	//校验和
	udph[18], udph[19] = 0x00, 0x00
	//计算校验值
	check := checkSum(append(udph, data...))
	udph[18], udph[19] = byte(check>>8&255), byte(check&255)

	return &udph
}
