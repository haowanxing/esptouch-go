package protocol

import (
	"github.com/haowanxing/esptouch-go/utils/byteutil"
	"net"
	"strconv"
	"strings"
)

type EsptouchGenerator struct {
	mGcBytes2 [][]byte
	mDcBytes2 [][]byte
}

func NewEsptouchGenerator(apSsid, apBssid, apPassword []byte, ipAddress net.IP) *EsptouchGenerator {
	var ipBytes = []byte{255, 255, 255, 255}
	for k, v := range strings.Split(ipAddress.String(), ".") {
		i, _ := strconv.Atoi(v)
		ipBytes[k] = byte(i)
	}
	gc := NewGuideCode()
	gcU81 := gc.GetU8s()
	mGcBytes2 := make([][]byte, len(gcU81))
	for i := 0; i < len(gcU81); i++ {
		mGcBytes2[i] = byteutil.GenSpecBytes(gcU81[i])
	}

	dc := NewDatumCode(apSsid, apBssid, apPassword, ipBytes)
	dcU81 := dc.GetU8s()
	mDcBytes2 := make([][]byte, len(dcU81))
	for i := 0; i < len(dcU81); i++ {
		mDcBytes2[i] = byteutil.GenSpecBytes(dcU81[i])
	}
	return &EsptouchGenerator{
		mGcBytes2: mGcBytes2,
		mDcBytes2: mDcBytes2,
	}
}

func (e *EsptouchGenerator) GetGCBytes2() [][]byte {
	return e.mGcBytes2
}
func (e *EsptouchGenerator) GetDCBytes2() [][]byte {
	return e.mDcBytes2
}
