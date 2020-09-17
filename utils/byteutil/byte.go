package byteutil

import (
	"fmt"
	"net"
)

func SplitUint8To2Bytes(char uint8) [2]byte {
	bt := [2]byte{0, 0}
	bt[0] = (char & 0xff) >> 4 //High
	bt[1] = char & 0x0f
	return bt
}

func Combine2BytesToOne(high, low byte) uint8 {
	return high<<4 | low
}

func CovertByte2Uint8(b byte) uint16 {
	return uint16(b & 0xff)
}
func Combine2bytesToU16(h, l byte) uint16 {
	highU8 := uint16(CovertByte2Uint8(h))
	lowU8 := CovertByte2Uint8(l)
	return highU8<<8 | uint16(lowU8)
}

func GenSpecBytes(length uint16) []byte {
	bt := make([]byte, 0)
	for i := uint16(0); i < length; i++ {
		bt = append(bt, '1')
	}
	return bt
}

func ParseBssid(bssidBytes []byte, offset, count int) string {
	bytes := bssidBytes[offset : offset+count]
	var sb string
	for _, v := range bytes {
		k := 0xff & v
		if k < 16 {
			sb += fmt.Sprintf("0%x", k)
		} else {
			sb += fmt.Sprintf("%x", k)
		}
	}
	return sb
}

func ParseInetAddr(inetAddrBytes []byte, offset, count int) net.IP {
	var sb string
	for i := 0; i < count; i++ {
		sb += fmt.Sprintf("%d", inetAddrBytes[offset+i]&0xff)
		if i != count-1 {
			sb += "."
		}
	}
	return net.ParseIP(sb)
}
