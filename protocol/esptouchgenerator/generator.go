package esptouchgenerator

import (
	"ESPTouch-Demo/protocol/datumcode"
	"ESPTouch-Demo/protocol/guidecode"
	"ESPTouch-Demo/utils/byteutil"
)

type EsptouchGenerator struct {
	mGcBytes2 [][]byte
	mDcBytes2 [][]byte
}

func NewEsptouchGenerator(apSsid, apBssid, apPassword, ipAddress []byte) *EsptouchGenerator {
	gc := guidecode.NewGuideCode()
	gcU81 := gc.GetU8s()
	mGcBytes2 := make([][]byte, len(gcU81))
	for i := 0; i < len(gcU81); i++ {
		mGcBytes2[i] = byteutil.GenSpecBytes(gcU81[i])
	}

	dc := datumcode.NewDatumCode(apSsid, apBssid, apPassword, ipAddress)
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
