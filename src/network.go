package main

import (
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

//GenIpfixPkt 生成ipfix报文
func GenIpfixPkt(srcIP, desIP, data string) {

}

//GenUDPPkt 生成udp报文
func GenUDPPkt(srcIP, desIP net.IP, buff string) *[]byte {
	//填充udp首部
	//udp伪首部
	udph := make([]byte, 20)
	//源ip地址
	udph[0], udph[1], udph[2], udph[3] = srcIP[12], srcIP[13], srcIP[14], srcIP[15]
	//目的ip地址
	udph[4], udph[5], udph[6], udph[7] = desIP[12], desIP[13], desIP[14], desIP[15]
	//协议类型
	udph[8], udph[9] = 0x00, 0x11
	//udp头长度
	udph[10], udph[11] = 0x00, byte(len(buff)+8)
	//下面开始就真正的udp头部
	//源端口号
	udph[12], udph[13] = 0x27, 0x10
	//目的端口号
	udph[14], udph[15] = 0x17, 0x70
	//udp头长度
	udph[16], udph[17] = 0x00, byte(len(buff)+8)
	//校验和
	udph[18], udph[19] = 0x00, 0x00
	//计算校验值
	check := checkSum(append(udph, buff...))
	udph[18], udph[19] = byte(check>>8&255), byte(check&255)

	return &udph
}

//decimal2Hex 将1个int转换为2个byte。exp：6000->0x17,0x70
func decimal2Hex(dec int) (byte,byte) {
	return byte(dec >> 8), byte(dec & ((1 << 8) - 1))
}

//GenIpv4Header 生成ipv4报文头
func GenIpv4Header(src, dst net.IP, buff string) *ipv4.Header {
	//填充ip首部
	iph := &ipv4.Header{
		Version: ipv4.Version,
		//IP头长一般是20
		Len: ipv4.HeaderLen,
		TOS: 0x00,
		//buff为数据
		TotalLen: ipv4.HeaderLen + len(buff),
		TTL:      64,
		Flags:    ipv4.DontFragment,
		FragOff:  0,
		Protocol: 17,
		Checksum: 0,
		Src:      src,
		Dst:      dst,
	}

	h, err := iph.Marshal()
	if err != nil {
		log.Fatalln(err)
	}
	//计算IP头部校验值
	iph.Checksum = int(checkSum(h))

	return iph
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
