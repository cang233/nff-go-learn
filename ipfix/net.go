package ipfix

import (
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

const GLOBAL_MIN_SEQUENCE = 1 //序列号起始值
const GLOBAL_MAX_SEQUENCE = 1<<16 // 序列号到最大时重置

var ipfix_sequence int //TODO 递增序列号,但由于要保证原子性,会影响性能，先跳过

//SendingIPFIX 发送IPFIX报文
func SendingIPFIX(srcIP, dstIP string, srcPort, dstPort int, ipfixData *IPFIX) {
	src := net.ParseIP(srcIP)
	dst := net.ParseIP(dstIP)

	// ipfixData.Hdr.PkgSequence = int32(ipfix_sequence)
	ipxBt := ipfixData.ToBytes()
	ipv4Hdr, _ := GenerateIPHeader(src, dst, len(ipxBt))
	udpHdr := GenerateUDPHeader(src, dst, srcPort, dstPort, ipxBt)
	send(ipv4Hdr,udpHdr,ipxBt)
}

func send(iph *ipv4.Header, udph []byte, buff []byte) {
	listener, err := net.ListenPacket("ip4:udp", iph.Src.String())
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	//listener 实现了net.PacketConn接口
	r, err := ipv4.NewRawConn(listener)
	if err != nil {
		log.Fatal(err)
	}
	//发送自己构造的UDP包
	if err = r.WriteTo(iph, append(udph[12:20], buff...), nil); err != nil {
		log.Fatal(err)
	}
}
