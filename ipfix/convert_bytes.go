package ipfix

func header2Bytes(hdr *Header) []byte {
	//gen ipfix header
	ipfixh := make([]byte, 20)
	//Version
	ipfixh[0], ipfixh[1] = Int16ToBytes(hdr.Version)
	//Count
	ipfixh[2], ipfixh[3] = Int16ToBytes(hdr.Count)
	//System up time
	ipfixh[4], ipfixh[5], ipfixh[6], ipfixh[7] = Int32ToBytes(hdr.SystemUpTime)
	// UNIX seconds
	ipfixh[8], ipfixh[9], ipfixh[10], ipfixh[11] = Int32ToBytes(hdr.UnixSeconds)
	// package Sequence
	ipfixh[12], ipfixh[13], ipfixh[14], ipfixh[15] = Int32ToBytes(hdr.PkgSequence)
	// source ID
	ipfixh[16], ipfixh[17], ipfixh[18], ipfixh[19] = Int32ToBytes(hdr.SourceID)

	return ipfixh
}

func flowset2Bytes(fset *FlowSet) []byte {
	flowSet := make([]byte, fset.Length)
	flowSet[0], flowSet[1] = Int16ToBytes(fset.ID)
	flowSet[2], flowSet[3] = Int16ToBytes(fset.Length)

	index := 4
	for _, v := range fset.Templates {
		flowSet[index], flowSet[index+1] = Int16ToBytes(v.ID)
		flowSet[index+2], flowSet[index+3] = Int16ToBytes(v.FieldCount)
		index += 4
		for _, f := range v.Fields {
			flowSet[index], flowSet[index+1] = Int16ToBytes(f.Type)
			flowSet[index+2], flowSet[index+3] = Int16ToBytes(f.Length)
			index += 4
		}
	}
	return flowSet
}

func optionFlowset2Bytes(opFset *OptionFlowSet) []byte {
	optionFlowSet := make([]byte,opFset.Length)
	optionFlowSet[0],optionFlowSet[1] = Int16ToBytes(opFset.ID)
	optionFlowSet[2],optionFlowSet[3] = Int16ToBytes(opFset.Length)

	index := 4
	for _,v := range opFset.Templates{
		optionFlowSet[index+0],optionFlowSet[index+1] = Int16ToBytes(v.ID)
		optionFlowSet[index+2],optionFlowSet[index+3] = Int16ToBytes(v.OptionScopeLength)
		optionFlowSet[index+4],optionFlowSet[index+5] = Int16ToBytes(v.OptionLength)
		index +=6
		for _, f := range v.ScopeFields {
			optionFlowSet[index], optionFlowSet[index+1] = Int16ToBytes(f.Type)
			optionFlowSet[index+2], optionFlowSet[index+3] = Int16ToBytes(f.Length)
			index += 4
		}
		for _, f := range v.OptionFields {
			optionFlowSet[index], optionFlowSet[index+1] = Int16ToBytes(f.Type)
			optionFlowSet[index+2], optionFlowSet[index+3] = Int16ToBytes(f.Length)
			index += 4
		}
	}
	return optionFlowSet
}

func dataFlowset2Bytes(dataFset *DataFlowSet) []byte {
	dataFlowSet := make([]byte,dataFset.Length)
	index := 0
	for _,v := range dataFset.Templates{
		dataFlowSet[index],dataFlowSet[index+1] = Int16ToBytes(v.ID)
		dataFlowSet[index+2],dataFlowSet[index+3] = Int16ToBytes(v.FieldCount)
		index +=4
		for _,f := range v.Fields{
			for i:=0;i<int(f.Length);i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += int(f.Length)
		}
	}
	for _,v := range dataFset.OptionTemplates{
		dataFlowSet[index],dataFlowSet[index+1] = Int16ToBytes(v.ID)
		dataFlowSet[index+2],dataFlowSet[index+3] = Int16ToBytes(v.OptionScopeLength+v.OptionLength)
		index +=4
		for _,f := range v.ScopeFields{
			for i:=0;i<int(f.Length);i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += int(f.Length)
		}
		for _,f := range v.OptionFields{
			for i:=0;i<int(f.Length);i++ {
				dataFlowSet[index+i] = f.Value[i]
			}
			index += int(f.Length)
		}
	}
	return dataFlowSet
}