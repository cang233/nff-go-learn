package ipfix

const NETFLOW_VERSION = 9 //IPFIX报文使用netflow v9报文格式

const (
	FLOW_SET_ID        = 0 // 流记录模板组成的FlowSet使用ID 0
	OPTION_FLOW_SET_ID = 1 // 选项模板组成的FlowSet使用ID 1
)

const TEMPLATE_ID_BEGIN = 256  //template id从此往后递增
const IPFIX_HEADER_LENGTH = 20 // ipfix header length

//IPFIX IPFIX报文
type IPFIX struct {
	Hdr      *Header
	Fset     *FlowSet
	OpFset   *OptionFlowSet
	DataFset *DataFlowSet

	Length int16 // 整个ipfix报文的长度,byte
}

//Header ipfix报文头信息
type Header struct {
	Version      int16 //ipfix报文格式,netflow v9版本
	Count        int16 //报文中携带的flow数量,个
	SystemUpTime int32 //设备运行的时间，以ms为单位
	UnixSeconds  int32 //从UTC时间1700 0时至现在的秒数
	PkgSequence  int32 //报文序列号，依次累加
	SourceID     int32 //用来标识Exporter所在的观察域，收集器可以根据数据包的源ip和sourceID来区分同一个exporter输出的不同输出流
}

//FlowSet 模板流组合
type FlowSet struct {
	ID        int16      //流记录模板组成的FlowSet使用ID 0
	Length    int16      //FlowSet的总长度(包括头+数据部分),byte
	Templates []Template //模板集合
}

//Template 模板
type Template struct {
	ID         int16   //用于数据与模板的对应。从256开始
	FieldCount int16   //Template record的字段数
	Fields     []Field //字段类型，用数字表示
}

//OptionFlowSet 选项记录模板流组合
type OptionFlowSet struct {
	ID        int16            //选项模板组成的FlowSet使用ID 1
	Length    int16            //FlowSet的长度，包括Padding长度（必有1个padding）,byte
	Templates []OptionTemplate //option模板集合
}

//OptionTemplate option的模板
type OptionTemplate struct {
	ID                int16   //用于数据与模板的对应，大于255
	OptionScopeLength int16   //Scope字段的字节数,byte
	OptionLength      int16   //Option字段的字节数,byte
	ScopeFields       []Field //：IPFIX进程相关数据引用的Scope字段类型。0x1：系统；0x2：接口；0x3：线卡；0x4：IPFIX cache；0x5：Template
	OptionFields      []Field //Option数据类型，使用的数值同流模板中介绍的Filed Type数值
}

//DataFlowSet 数据流组合。
//可能有padding：用于使FlowSet的长度按照32位圆整。
type DataFlowSet struct {
	Templates       []Template
	OptionTemplates []OptionTemplate
	Length          int16 //Data flowset的总长度，byte
}

//Field 流字段
type Field struct {
	Type   int16  //Field Type,字段类型，用数字表示,rfc:https://www.iana.org/assignments/ipfix/ipfix.xml#ipfix-information-elements
	Length int16  //Field Length,单位:byte
	Value  []byte //Field Value
}

//Init 初始化一些ipfix报文中的字段
func Init() *IPFIX {
	return &IPFIX{
		Hdr: &Header{
			Version:  NETFLOW_VERSION,
			Count:    3,
			SourceID: 0,
		},
		Fset: &FlowSet{
			ID:     FLOW_SET_ID,
			Length: 4,
		},
		OpFset: &OptionFlowSet{
			ID:     OPTION_FLOW_SET_ID,
			Length: 4 + 2, //默认4个,必有1个padding
		},
	}
}

//AddField 添加flowset的template
func (ipx *IPFIX) AddField(template Template) {
	//添加模板
	ipx.Fset.Templates = append(ipx.Fset.Templates, template)
	//将数据保存到data flowset中
	ipx.DataFset.Templates = append(ipx.DataFset.Templates, template)
	//将模板长度累加到flowset中
	for _, f := range template.Fields {
		ipx.Fset.Length += f.Length + 2
		ipx.DataFset.Length += f.Length
	}
	ipx.Fset.Length += 4     // +id,length
	ipx.DataFset.Length += 4 // +id,length
}

//AddOptionField 添加option flowset的template
func (ipx *IPFIX) AddOptionField(optionTemplate OptionTemplate) {
	ipx.OpFset.Templates = append(ipx.OpFset.Templates, optionTemplate)
	ipx.OpFset.Length += optionTemplate.OptionScopeLength + optionTemplate.OptionLength + 6 //+id,scope length,option length
	ipx.DataFset.OptionTemplates = append(ipx.DataFset.OptionTemplates, optionTemplate)
	ipx.DataFset.Length += (optionTemplate.OptionScopeLength+optionTemplate.OptionLength)>>2 + 4 //+id,length
}

//Fill 完善剩余字段，生成最终IPFIX报文
func (ipx *IPFIX) Fill() {
	ipx.Length += IPFIX_HEADER_LENGTH + ipx.Fset.Length + ipx.OpFset.Length + ipx.DataFset.Length
}

//ToBytes 将IPFIX报文转化为[]byte
func (ipx *IPFIX) ToBytes() []byte {
	var ipxBs []byte
	ipxBs = append(ipxBs, header2Bytes(ipx.Hdr)...)
	ipxBs = append(ipxBs, flowset2Bytes(ipx.Fset)...)
	ipxBs = append(ipxBs, optionFlowset2Bytes(ipx.OpFset)...)
	ipxBs = append(ipxBs, dataFlowset2Bytes(ipx.DataFset)...)
	return ipxBs
}
