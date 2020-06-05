package byteutil

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
