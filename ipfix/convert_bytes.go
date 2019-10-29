package ipfix

import (
	"log"
)

func header2Bytes(hdr *Header) *[]byte {
	//gen ipfix header
	ipfixh := make([]byte, IPFIX_HEADER_LENGTH)
	//Version
	ipfixh[0], ipfixh[1] = Int16ToBytes(hdr.Version)
	//Count
	ipfixh[2], ipfixh[3] = Int16ToBytes(hdr.Count)
	//System up time
	ipfixh[4], ipfixh[5], ipfixh[6], ipfixh[7] = Int32ToBytes(hdr.SystemUpTime)
	// UNIX seconds
	ipfixh[8], ipfixh[9], ipfixh[10], ipfixh[11] = Int64ToBytes(hdr.UnixSeconds)
	// package Sequence
	ipfixh[12], ipfixh[13], ipfixh[14], ipfixh[15] = Int32ToBytes(hdr.PkgSequence)
	// source ID
	ipfixh[16], ipfixh[17], ipfixh[18], ipfixh[19] = Int32ToBytes(hdr.SourceID)

	return &ipfixh
}

func flowset2Bytes(fset *FlowSet) *[]byte {
	fsetLen := 4
	for _, t := range fset.Templates {
		fsetLen += 4 + len(t.Fields)*4
	}
	log.Println("FlowSet pkt len :", fsetLen)

	flowSet := make([]byte, fsetLen)
	flowSet[0], flowSet[1] = Int16ToBytes(fset.ID)
	flowSet[2], flowSet[3] = IntToBytes(fsetLen)

	index := 4
	for _, v := range fset.Templates {
		flowSet[index], flowSet[index+1] = Int16ToBytes(v.ID)
		flowSet[index+2], flowSet[index+3] = Int16ToBytes(int16(len(v.Fields)))
		index += 4
		for _, f := range v.Fields {
			flowSet[index], flowSet[index+1] = Int16ToBytes(f.Type)
			flowSet[index+2], flowSet[index+3] = Int16ToBytes(int16(len(f.Value)))
			index += 4
		}
	}
	log.Printf("FlowSet=%+v\n", flowSet)
	return &flowSet
}

func optionFlowset2Bytes(opFset *OptionFlowSet) *[]byte {
	totalLen := 4
	for _, t := range opFset.Templates {
		totalLen += 4 * len(t.ScopeFields)
		totalLen += 4 * len(t.OptionFields)
		totalLen += 6
	}
	if totalLen %4==2{
		totalLen += 2
	}
	log.Println("OptionFlowSet pkt len :", totalLen)

	optionFlowSet := make([]byte, totalLen)
	optionFlowSet[0], optionFlowSet[1] = Int16ToBytes(opFset.ID)
	optionFlowSet[2], optionFlowSet[3] = IntToBytes(totalLen)

	index := 4
	for _, v := range opFset.Templates {
		optionFlowSet[index+0], optionFlowSet[index+1] = Int16ToBytes(v.ID)
		optionFlowSet[index+2], optionFlowSet[index+3] = IntToBytes(len(v.ScopeFields) * 4)
		optionFlowSet[index+4], optionFlowSet[index+5] = IntToBytes(len(v.OptionFields) * 4)
		index += 6
		for _, f := range v.ScopeFields {
			optionFlowSet[index], optionFlowSet[index+1] = Int16ToBytes(f.Type)
			optionFlowSet[index+2], optionFlowSet[index+3] = IntToBytes(len(f.Value))
			index += 4
		}
		for _, f := range v.OptionFields {
			optionFlowSet[index], optionFlowSet[index+1] = Int16ToBytes(f.Type)
			optionFlowSet[index+2], optionFlowSet[index+3] = IntToBytes(len(f.Value))
			index += 4
		}
	}
	log.Printf("OptionFlowSet=%+v\n", optionFlowSet)
	return &optionFlowSet
}

func dataFlowset2Bytes(dataFset *DataFlowSet) *[]byte {
	dataFsLen := 0
	for _, t := range dataFset.Templates {
		tpLen := 4
		for _, f := range t.Fields {
			tpLen += len(f.Value)
		}
		//dataFlow部分，一个template就对应一个flowset,因此每个flowset都需要填充
		if tpLen%4 == 2 {
			tpLen += 2
		}
		dataFsLen += tpLen
	}
	for _, t := range dataFset.OptionTemplates {
		otpLen := 4
		for _, f := range t.ScopeFields {
			otpLen += len(f.Value)
		}
		for _, f := range t.OptionFields {
			otpLen += len(f.Value)
		}
		//dataFlow部分，一个template就对应一个flowset,因此每个flowset都需要填充
		if otpLen%4 == 2 {
			otpLen += 2
		}
		dataFsLen += otpLen
	}

	log.Println("DataFlowSet pkt len :", dataFsLen)

	dataFlowSet := make([]byte, dataFsLen)
	index := 0
	for _, v := range dataFset.Templates {
		flen := 4
		needPadding := false
		for _, f := range v.Fields {
			flen += len(f.Value)
		}
		if flen%4 == 2 {
			flen += 2
			needPadding = true
		}
		dataFlowSet[index], dataFlowSet[index+1] = Int16ToBytes(v.ID)
		dataFlowSet[index+2], dataFlowSet[index+3] = Int16ToBytes(int16(flen))
		index += 4
		for _, f := range v.Fields {
			for i := 0; i < len(f.Value); i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += len(f.Value)
		}
		if needPadding {
			index += 2
		}
	}
	for _, v := range dataFset.OptionTemplates {
		opLen := 4
		needPadding := false
		for _, os := range v.ScopeFields {
			opLen += len(os.Value)
		}
		for _, os := range v.OptionFields {
			opLen += len(os.Value)
		}
		if opLen%4 == 2 {
			needPadding = true
			opLen += 2
		}
		dataFlowSet[index], dataFlowSet[index+1] = Int16ToBytes(v.ID)
		dataFlowSet[index+2], dataFlowSet[index+3] = IntToBytes(opLen)
		index += 4
		for _, f := range v.ScopeFields {
			for i := 0; i < len(f.Value); i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += len(f.Value)
		}
		for _, f := range v.OptionFields {
			for i := 0; i < len(f.Value); i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += len(f.Value)
		}
		if needPadding {
			index += 2
		}
	}

	log.Printf("DataFlowSet=%+v\n", dataFlowSet)
	return &dataFlowSet
}
