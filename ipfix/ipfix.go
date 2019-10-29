package ipfix

import (
	"log"
	"time"
)

const NETFLOW_VERSION = 9 //IPFIX报文使用netflow v9报文格式

const (
	FLOW_SET_ID        int16 = 0 // 流记录模板组成的FlowSet使用ID 0
	OPTION_FLOW_SET_ID int16 = 1 // 选项模板组成的FlowSet使用ID 1
)

const TEMPLATE_ID_BEGIN int16 = 256 //template id从此往后递增
const IPFIX_HEADER_LENGTH = 20      // ipfix header length

//IPFIX IPFIX报文
type IPFIX struct {
	Hdr      *Header
	Fset     *FlowSet
	OpFset   *OptionFlowSet
	DataFset *DataFlowSet
}

//Header ipfix报文头信息
type Header struct {
	Version      int16 //ipfix报文格式,netflow v9版本
	Count        int16 //报文中携带的flow数量,个
	SystemUpTime int32 //设备运行的时间，以ms为单位
	UnixSeconds  int64 //从UTC时间1700 0时至现在的秒数
	PkgSequence  int32 //报文序列号，依次累加
	SourceID     int32 //用来标识Exporter所在的观察域，收集器可以根据数据包的源ip和sourceID来区分同一个exporter输出的不同输出流
}

//FlowSet 模板流组合
type FlowSet struct {
	ID        int16      //流记录模板组成的FlowSet使用ID 0
	Templates []Template //模板集合
}

//Template 模板
type Template struct {
	ID     int16   //用于数据与模板的对应。从256开始
	Fields []Field //字段类型，用数字表示
}

//OptionFlowSet 选项记录模板流组合
type OptionFlowSet struct {
	ID        int16            //选项模板组成的FlowSet使用ID 1
	Templates []OptionTemplate //option模板集合
}

//OptionTemplate option的模板
type OptionTemplate struct {
	ID           int16   //用于数据与模板的对应，大于255
	ScopeFields  []Field //：IPFIX进程相关数据引用的Scope字段类型。0x1：系统；0x2：接口；0x3：线卡；0x4：IPFIX cache；0x5：Template
	OptionFields []Field //Option数据类型，使用的数值同流模板中介绍的Filed Type数值
}

//DataFlowSet 数据流组合。
//可能有padding：用于使FlowSet的长度按照32位圆整。
type DataFlowSet struct {
	Templates       []Template
	OptionTemplates []OptionTemplate
}

//Field 流字段
type Field struct {
	Type  int16  //Field Type,字段类型，用数字表示,rfc:https://www.iana.org/assignments/ipfix/ipfix.xml#ipfix-information-elements
	Value []byte //Field Value
}

//Init 新建IPFIX对象，初始化一些ipfix报文中的字段并返回
func Init() *IPFIX {
	return &IPFIX{
		Hdr: &Header{
			Version:     NETFLOW_VERSION,
			Count:       2,
			SourceID:    0,
			UnixSeconds: time.Now().Unix(),
		},
		Fset: &FlowSet{
			ID: FLOW_SET_ID,
		},
		OpFset: &OptionFlowSet{
			ID: OPTION_FLOW_SET_ID,
		},
		DataFset: &DataFlowSet{},
	}
}

//AddTemplate 添加flowset的template
func (ipx *IPFIX) AddTemplate(template Template) {
	//添加模板
	ipx.Fset.Templates = append(ipx.Fset.Templates, template)
	//将数据保存到data flowset中
	ipx.DataFset.Templates = append(ipx.DataFset.Templates, template)
	//将模板长度累加到flowset中
	ipx.Hdr.Count++ //+flow数量
}

//AddOptionTemplate 添加option flowset的template
func (ipx *IPFIX) AddOptionTemplate(optionTemplate OptionTemplate) {
	ipx.OpFset.Templates = append(ipx.OpFset.Templates, optionTemplate)
	ipx.DataFset.OptionTemplates = append(ipx.DataFset.OptionTemplates, optionTemplate)
	ipx.Hdr.Count++ //+flow数量
}

//Fill 完善剩余字段，生成最终IPFIX对象 //TODO
func (ipx *IPFIX) Fill() {
}

//ToBytes 将IPFIX报文转化为[]byte
func (ipx *IPFIX) ToBytes() []byte {
	var ipxBs []byte
	ipxBs = append(ipxBs, *header2Bytes(ipx.Hdr)...)
	ipxBs = append(ipxBs, *flowset2Bytes(ipx.Fset)...)
	ipxBs = append(ipxBs, *optionFlowset2Bytes(ipx.OpFset)...)
	ipxBs = append(ipxBs, *dataFlowset2Bytes(ipx.DataFset)...)

	log.Println("IPFIX pkt len :", len(ipxBs))
	log.Printf("IPFIX pkt:%+v\n", ipxBs)
	return ipxBs
}
