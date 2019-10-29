package test

import (
	"log"
	"net"
	"nff-go-learn/ipfix"
	"testing"

	"golang.org/x/net/ipv4"
)

func TestGenIpfix(t *testing.T) {
	ipx := ipfix.Init()
	ipx.Hdr.SystemUpTime = 15
	//每次更新需要全更新ipx报文所有的再重新生成一遍二进制？？？ //TODO
	seq := 0

	ipx.Fill()

	sip := "10.10.26.17"
	dip := "10.10.28.47"
	src := net.ParseIP(sip)
	dst := net.ParseIP(dip)
	srcp := 10000
	dstp := 2055

	templateID := ipfix.TEMPLATE_ID_BEGIN + 1
	tmplte := ipfix.Template{
		ID: templateID,
		Fields: []ipfix.Field{
			ipfix.Field{
				Type:  8,
				Value: []byte{src[12], src[13], src[14], src[15]},
			},
			ipfix.Field{
				Type:  12,
				Value: []byte{dst[12], dst[13], dst[14], dst[15]},
			},
		},
	}
	ipx.AddTemplate(tmplte)

	templateID++
	tmplte = ipfix.Template{
		ID: templateID,
		Fields: []ipfix.Field{
			ipfix.Field{
				Type:  12,
				Value: []byte{dst[12], dst[13], dst[14], dst[15]},
			},
		},
	}
	ipx.AddTemplate(tmplte)

	for seq = 0; seq < 10000; seq++ {
		ipx.Hdr.PkgSequence = int32(seq)
		ipxBt := ipx.ToBytes()
		ipv4Hdr, _ := ipfix.GenerateIPHeader(src, dst, len(ipxBt))
		udpHdr := ipfix.GenerateUDPHeader(src, dst, srcp, dstp, ipxBt)
		sendPacket(ipv4Hdr, udpHdr, ipxBt)
	}

}

func TestSendEmpty(t *testing.T) {
	seq := 0

	sip := "10.10.26.17"
	dip := "10.10.28.47"
	src := net.ParseIP(sip)
	dst := net.ParseIP(dip)
	srcp := 10000
	dstp := 2055

	data := []byte{}

	for seq = 0; seq < 100; seq++ {
		// ipx.Hdr.PkgSequence = int32(seq)
		ipv4Hdr, _ := ipfix.GenerateIPHeader(src, dst, len(data))
		udpHdr := ipfix.GenerateUDPHeader(src, dst, srcp, dstp, data)
		sendPacket(ipv4Hdr, udpHdr, data)
	}
}

func sendPacket(iph *ipv4.Header, udph []byte, buff []byte) {
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
