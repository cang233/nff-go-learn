package test

import (
	"net"
	"nff-go-learn/ipfix"
	"testing"
	"time"
)

func TestGenIpfix(t *testing.T) {
	ipx := &ipfix.IPFIX{}
	ipx.Hdr.UnixSeconds = int32(time.Now().Unix())
	//每次更新需要全更新ipx报文所有的再重新生成一遍二进制？？？ //TODO
	seq := 0
	templateID := ipfix.TEMPLATE_ID_BEGIN

	ipx.Fill()

	src := net.ParseIP("10.10.26.17")
	dst := net.ParseIP("10.10.26.47")
	srcp := 50321
	dstp := 2055

	for seq = 0; seq < 20; seq++ {
		ipx.Hdr.PkgSequence = int32(seq)
		ipxBt := ipx.ToBytes()
		ipv4Hdr,_ := ipfix.GenerateIPHeader(src,dst,ipxBt)
		udpHdr := ipfix.GenerateUDPHeader(src,dst,srcp,dstp,ipxBt)
	}

}
