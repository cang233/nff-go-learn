package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"
	"github.com/intel-go/nff-go/types"
)

func main() {
	flag.UintVar(&size, "s", 64, "size of packets,(bit)")
	flag.UintVar(&speed, "v", 1, "speed of generation,(GB)")
	flag.Parse()

	pkts := uint64(speed *1000 * 1000 * 1000 / 8 / (size + 20))
	size = size - types.EtherLen - types.IPv4MinLen - types.TCPMinLen - 4 /* Ethernet checksum length*/

	statsServerAddres := &net.TCPAddr{
		Port: 8080,
	}
	flow.SystemInit(&flow.Config{
//		CPUList:          "0-23",
		StatsHTTPAddress: statsServerAddres,
	})
	// rev, err := flow.SetReceiver(0)

	// generator默认1GB一个,大于1GB的分多个。
	var rev *flow.Flow
	var err error
	if(speed<=1){
		rev, _, err = flow.SetFastGenerator(generatePacket, pkts, *new(ctx))
		flow.CheckFatal(err)
	}else{
		var flows []*flow.Flow
		var i uint
		for i = 0;i<speed;i++{
			fl,_,err:=flow.SetFastGenerator(generatePacket, pkts/uint64(speed), *new(ctx))
			flow.CheckFatal(err)
			flows = append(flows,fl) 
		}
		rev,err = flow.SetMerger(flows...)
		flow.CheckFatal(err)
	}

	flow.CheckFatal(err)
	flow.CheckFatal(flow.SetHandler(rev, parseOnlyNoPrint, *new(ctx)))
	flow.CheckFatal(flow.SetStopper(rev))

	flow.CheckFatal(flow.SystemStart())
}

func generater(p *packet.Packet, context flow.UserContext) {
	// Total packet size will be 14+20+20+70+4(crc)=128 bytes
	//step1: 设置种子数
	rand.Seed(time.Now().UnixNano())
	//step2：获取随机数
	flag := rand.Intn(10) //[0,100)
	switch flag {
	case 0:
		packet.InitEmptyARPPacket(p)
	case 1:
		packet.InitEmptyIPv4ICMPPacket(p, 70)
	case 2:
		packet.InitEmptyIPv4Packet(p, 70)
	case 3:
		packet.InitEmptyIPv4TCPPacket(p, 70)
	case 4:
		packet.InitEmptyIPv4UDPPacket(p, 70)
	case 5:
		packet.InitEmptyIPv6ICMPPacket(p, 70)
	case 6:
		packet.InitEmptyIPv6Packet(p, 70)
	case 7:
		packet.InitEmptyIPv6TCPPacket(p, 70)
	case 8:
		packet.InitEmptyIPv6UDPPacket(p, 70)
	default:
		packet.InitEmptyPacket(p, 70)
	}
}

func handlerF1(p *packet.Packet, u flow.UserContext) {
	p.ParseL3()
	if ipv4 := p.GetIPv4(); ipv4 != nil {
		fmt.Println("Package L3 type: IPV4")
		fmt.Println(ipv4.String())
		// fmt.Println(ipv4.SrcPort)
		// fmt.Println(ipv4.DstPort)
		// fmt.Println(ipv4.SentSeq)
		// fmt.Println(ipv4.RecvAck)
		// fmt.Println(ipv4.DataOff)
		// fmt.Println(ipv4.TCPFlags)
		// fmt.Println(ipv4.RxWin)
		// fmt.Println(ipv4.Cksum)
		// fmt.Println(ipv4.TCPUrp)

		p.ParseL4ForIPv4()
		if tcp := p.GetTCPForIPv4(); tcp != nil {
			fmt.Print("Package L4 type: TCP")
			fmt.Println(tcp.String())
		} else if udp := p.GetUDPForIPv4(); udp != nil {
			fmt.Print("Package L4 type: UDP")
			fmt.Println(udp.String())
		} else if icmp := p.GetICMPForIPv4(); icmp != nil {
			fmt.Print("Package L4 type: ICMP")
			fmt.Println(icmp.String())
		} else {
			fmt.Println("Unknow L4 type...")
		}

	} else if ipv6 := p.GetIPv6(); ipv6 != nil {
		fmt.Println("Package L3 type: IPV6")
		fmt.Println(ipv6.String())

		p.ParseL4ForIPv6()
		if tcp := p.GetTCPForIPv6(); tcp != nil {
			fmt.Print("Package L4 type: TCP")
			fmt.Println(tcp.String())
		} else if udp := p.GetUDPForIPv6(); udp != nil {
			fmt.Print("Package L4 type: UDP")
			fmt.Println(udp.String())
		} else if icmp := p.GetICMPForIPv6(); icmp != nil {
			fmt.Print("Package L4 type: ICMP")
			fmt.Println(icmp.String())
		} else {
			fmt.Println("Unknow L4 type...")
		}
	} else {
		fmt.Println("Unknow package protocol------")
		fmt.Println("GetEtherType : ", p.GetEtherType())
		fmt.Println("GetPacketLen : ", p.GetPacketLen())
		fmt.Println("GetRawPacketBytes : ", p.GetRawPacketBytes())
		fmt.Println("GetPacketTimestamp : ", p.GetPacketTimestamp())
		fmt.Println("GetPacketOffloadFlags : ", p.GetPacketOffloadFlags())
		fmt.Println("GetPacketSegmentLen : ", p.GetPacketSegmentLen())
	}
}

//parseOnlyNoPrint 虽然封装了AllKnown的爬取方式，但是看源码好像也没省多少代码和性能。
// 只解析，没有打印
func parseOnlyNoPrint(p *packet.Packet, context flow.UserContext) {
	ipv4, ipv6, arp := p.ParseAllKnownL3()
	if ipv4 != nil {
		tcp, udp, icmp := p.ParseAllKnownL4ForIPv4()
		if tcp != nil {
			sum := tcp.String()
			sum += strconv.Itoa(int(tcp.SrcPort))
			sum += strconv.Itoa(int(tcp.DstPort))
			sum += strconv.Itoa(int(ipv4.TypeOfService))
			sum += ipv4.SrcAddr.String()
			sum += ipv4.DstAddr.String()
		} else if udp != nil {
			sum := udp.String()
			sum += strconv.Itoa(int(udp.SrcPort))
			sum += strconv.Itoa(int(udp.DstPort))
			sum += strconv.Itoa(int(ipv4.TypeOfService))
			sum += ipv4.SrcAddr.String()
			sum += ipv4.DstAddr.String()
		} else if icmp != nil {
			sum := icmp.String()
			sum += strconv.Itoa(int(icmp.Type))
			sum += strconv.Itoa(int(icmp.Code))
			sum += strconv.Itoa(int(ipv4.TypeOfService))
			sum += ipv4.SrcAddr.String()
			sum += ipv4.DstAddr.String()
		}
	} else if ipv6 != nil {
		tcp, udp, icmp := p.ParseAllKnownL4ForIPv6()
		if tcp != nil {
			sum := tcp.String()
			sum += strconv.Itoa(int(tcp.SrcPort))
			sum += strconv.Itoa(int(tcp.DstPort))
			sum += ipv6.SrcAddr.String()
			sum += ipv6.DstAddr.String()
		} else if udp != nil {
			sum := udp.String()
			sum += strconv.Itoa(int(udp.SrcPort))
			sum += strconv.Itoa(int(udp.DstPort))
			sum += ipv6.SrcAddr.String()
			sum += ipv6.DstAddr.String()
		} else if icmp != nil {
			sum := icmp.String()
			sum += strconv.Itoa(int(icmp.Type))
			sum += strconv.Itoa(int(icmp.Code))
			sum += ipv6.SrcAddr.String()
			sum += ipv6.DstAddr.String()
		}
	} else if arp != nil {
		sum := arp.String()
		sum += arp.SHA.String()
		sum += arp.THA.String()
		sum += strconv.Itoa(int(arp.HType))
		sum += strconv.Itoa(int(arp.PType))
		sum += strconv.Itoa(int(arp.HLen))
		sum += strconv.Itoa(int(arp.PLen))
	} else {
		sum := "Error,"
		sum += "Not recognised package type."
	}

}

var count uint32 = 0
var lastTime time.Time

func counter(p *packet.Packet, context flow.UserContext) {
	if count == 0 {
		lastTime = time.Now()
	}
	duration := time.Now().Sub(lastTime).Seconds()
	count += 10000
	fmt.Println("--------------------------------------------")
	fmt.Printf("Total received %d packages at %s", count, time.Now().UTC().String())
	fmt.Println("Average process Speed: ", 10000/duration, " pps")
	fmt.Println("--------------------------------------------")

}

var size uint
var speed uint

type ctx struct {
	r uint16
}

func (c ctx) Copy() interface{} {
	n := new(ctx)
	n.r = 20
	return n
}

func (c ctx) Delete() {
}

// Function to use in generator
func generatePacket(pkt *packet.Packet, context flow.UserContext) {
	ctx1 := context.(*ctx)
	r := ctx1.r
	packet.InitEmptyIPv4TCPPacket(pkt, size)
	ipv4 := pkt.GetIPv4()
	tcp := pkt.GetTCPForIPv4()

	ipv4.DstAddr = packet.SwapBytesIPv4Addr(types.IPv4Address(r))
	ipv4.SrcAddr = packet.SwapBytesIPv4Addr(types.IPv4Address(r + 15))

	tcp.DstPort = packet.SwapBytesUint16(r + 25)
	tcp.SrcPort = packet.SwapBytesUint16(r + 35)

	ctx1.r++
	if ctx1.r > 259 {
		ctx1.r = 20
	}
}
