package test

import (
	"nff-go-learn/ipfix"
	"testing"
	"time"
)

func TestGenIpfix(t *testing.T) {
	ipx := &ipfix.IPFIX{}
	ipx.Hdr.UnixSeconds = int32(time.Now().Unix())
}
