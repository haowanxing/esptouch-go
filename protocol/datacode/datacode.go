package datacode

import "ESPTouch-Demo/utils/byteutil"
import "ESPTouch-Demo/utils"

const (
	DATA_CODE_LEN = 6
	INDEX_MAX     = 127
)

type DataCode struct {
	mSeqHeader byte
	mDataHigh  byte
	mDataLow   byte
	mCrcHigh   byte
	mCrcLow    byte
}

func NewDataCode(u8 uint16, index int) *DataCode {
	dataBytes := byteutil.SplitUint8To2Bytes(byte(u8))
	crc := utils.NewCRC8()
	crc.Update([]byte{byte(u8)}, 0, 1)
	crc.Update([]byte{byte(index)}, 0, 1)
	crcBytes := byteutil.SplitUint8To2Bytes(uint8(crc.GetValue()))
	return &DataCode{
		mSeqHeader: byte(index),
		mDataHigh:  dataBytes[0],
		mDataLow:   dataBytes[1],
		mCrcHigh:   crcBytes[0],
		mCrcLow:    crcBytes[1],
	}
}

func (d *DataCode) GetBytes() []byte {
	dataBytes := make([]byte, DATA_CODE_LEN)
	dataBytes[0] = 0x00
	dataBytes[1] = byteutil.Combine2BytesToOne(d.mCrcHigh, d.mDataHigh)
	dataBytes[2] = 0x01
	dataBytes[3] = d.mSeqHeader
	dataBytes[4] = 0x00
	dataBytes[5] = byteutil.Combine2BytesToOne(d.mCrcLow, d.mDataLow)
	return dataBytes
}
