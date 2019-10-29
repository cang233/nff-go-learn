package test

import (
	"strconv"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"golang.org/x/net/ipv4"
)

func TestParseIp(t *testing.T) {
	ip := net.ParseIP("10.10.26.24")
	fmt.Printf("%+v\n", ip[14])
}

func TestConvert(t *testing.T) {
	port := 6000
	b1 := byte(port >> 8)
	b2 := byte(port & ((1 << 8) - 1))
	fmt.Printf("%+v,%+v\n", b1, b2)
	fmt.Println(0x17, 0x70)

	bs :=IntToBytes(time.Now().Unix())
	fmt.Println(len(bs))
	fmt.Println(strconv.Itoa(int(time.Now().Unix())))
	fmt.Printf("%+v\n",bs)
}

func port2Hex(port int) (byte, byte) {
	return byte(port >> 8), byte(port & ((1 << 8) - 1))
}

func TestGenUdp(t *testing.T) {

	//目的IP
	dst := net.IPv4(10, 10, 28, 47)
	//源IP
	src := net.IPv4(10, 10, 26, 17)
	// buff := "this is dataadsfasdfasdfgsdffffffffffssdfgaqwefawerfvsergfqwefawefdasgqwrefvsaerfaqwefvsergswregfvaffffffhbgasfawefasdfqwefasdafasdfafd."

	//gen ipfix header
	ipfixh := make([]byte, 20)
	//Version
	ipfixh[0], ipfixh[1] = 0x00, 0x09
	//Count
	ipfixh[2], ipfixh[3] = port2Hex(4)

	//System up time
	ipfixh[4],ipfixh[5],ipfixh[6],ipfixh[7] = 0x00,0x00,0x00,0x12
	unixTime := IntToBytes(time.Now().Unix())
	// UNIX seconds
	fmt.Println("Unix len:",len(unixTime))
	ipfixh[8],ipfixh[9],ipfixh[10],ipfixh[11] = unixTime[4],unixTime[5],unixTime[6],unixTime[7]
	// package Sequence
	ipfixh[12],ipfixh[13],ipfixh[14],ipfixh[15] = 0x00,0x00,0x00,0x01
	// source ID
	ipfixh[16],ipfixh[17],ipfixh[18],ipfixh[19] = 0x00,0x00,0x00,0x00

	//FlowFlowSet
	flowSet := make([]byte,24)
	//flowSet ID
	flowSet[0],flowSet[1] = 0x00,0x00
	//length
	flowSet[2],flowSet[3] = port2Hex(24)
	//template ID 	
	flowSet[4],flowSet[5] = port2Hex(257)
	//Field count
	flowSet[6],flowSet[7] = 0x00,0x02
	//Field 1 Type
	flowSet[8],flowSet[9] = port2Hex(8)
	//Field 1 type Length,byte
	flowSet[10],flowSet[11] = 0x00,0x04
	//Field 2 Type
	flowSet[12],flowSet[13] = port2Hex(12)
	//Field 2 type Length,byte
	flowSet[14],flowSet[15] = 0x00,0x04
	//template2 ID
	flowSet[16],flowSet[17] = port2Hex(258)
	//Field count
	flowSet[18],flowSet[19] = 0x00,0x01
	//Field 3 Type
	flowSet[20],flowSet[21] = port2Hex(12)
	//Field 4 type Length,byte
	flowSet[22],flowSet[23] = 0x00,0x04


	//option flow set
	optionFlowSet := make([]byte,4)
	//optional flow set id
	optionFlowSet[0],optionFlowSet[1] = 0x00,0x01
	// optional flow set length
	optionFlowSet[2],optionFlowSet[3] = 0x00,0x04
	

	//Data Flow Set
	dataFlowSet := make([]byte,20)
	//data flow set id
	dataFlowSet[0],dataFlowSet[1] = port2Hex(257)
	//data flow set length
	dataFlowSet[2],dataFlowSet[3] = port2Hex(12)
	//Field 1 value
	dataFlowSet[4],dataFlowSet[5],dataFlowSet[6],dataFlowSet[7] = src[12], src[13], src[14], src[15]
	//Field 2 value
	dataFlowSet[8],dataFlowSet[9],dataFlowSet[10],dataFlowSet[11] = dst[12], dst[13], dst[14], dst[15]
	//data flow set id
	dataFlowSet[12],dataFlowSet[13] = port2Hex(258)
	//data flow set length
	dataFlowSet[14],dataFlowSet[15] = port2Hex(8)
	//Field 1 value
	dataFlowSet[16],dataFlowSet[17],dataFlowSet[18],dataFlowSet[19] = dst[12], dst[13], dst[14], dst[15]

	buff := append(ipfixh,append(flowSet,append(optionFlowSet,dataFlowSet...)...)...)

	// var buff []byte
	// buff = append(buff,ipfixh...)

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

	//填充udp首部
	//udp伪首部
	udph := make([]byte, 20)
	//源ip地址
	udph[0], udph[1], udph[2], udph[3] = src[12], src[13], src[14], src[15]
	//目的ip地址
	udph[4], udph[5], udph[6], udph[7] = dst[12], dst[13], dst[14], dst[15]
	//协议类型
	udph[8], udph[9] = 0x00, 0x11
	//udp头长度
	udph[10], udph[11] = 0x00, byte(len(buff)+8)
	//下面开始就真正的udp头部
	//源端口号
	udph[12], udph[13] = 0x27, 0x10
	//目的端口号，改成2055后wireshark就自动解析成cflow协议了。
	udph[14], udph[15] = port2Hex(2055)
	//udp头长度
	udph[16], udph[17] = 0x00, byte(len(buff)+8)
	//校验和
	udph[18], udph[19] = 0x00, 0x00
	//计算校验值
	check := checkSum(append(udph, buff...))
	udph[18], udph[19] = byte(check>>8&255), byte(check&255)
	
	listener, err := net.ListenPacket("ip4:udp", "10.10.26.17")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	//listener 实现了net.PacketConn接口
	r, err := ipv4.NewRawConn(listener)
	if err != nil {
		log.Fatal(err)
	}

	index := 0
	//发送自己构造的UDP包
	for {
		if err = r.WriteTo(iph, append(udph[12:20], buff...), nil); err != nil {
			log.Fatal(err)
		}
		fmt.Println("No.", index)
		index++
		if index > 50 {
			return
		}
	}
}
func IntToBytes(data int64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
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
