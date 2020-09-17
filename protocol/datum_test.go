package protocol

import (
	"encoding/binary"
	"testing"
)

func TestNewDatumCode(t *testing.T) {
	dc := NewDatumCode([]byte("Administrators"), []byte{0x00, 0x1f, 0x7a, 0x71, 0x93, 0xb0}, []byte("123qweasdzxc"), []byte{192, 168, 123, 196})
	t.Log(len(dc.GetBytes()), dc.GetBytes())
	t.Log(len(dc.GetU8s()), dc.GetU8s())
}

func TestDatumCode_GetBytes(t *testing.T) {
	a := binary.LittleEndian.Uint64([]byte{227, 4, 139, 108, 114, 1, 0, 0})
	t.Log(a)
}
