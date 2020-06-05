package utils

import "testing"

func TestCRC8_GetValue(t *testing.T) {
	//t.Logf("%x",Crc8([]byte{0x1a, 0x0}))
	crc := NewCRC8()
	crc.Update([]byte{0x1a, 0x1a}, 0, 2)
	t.Logf("%.2x", crc.GetValue())
}
