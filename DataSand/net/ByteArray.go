package net

import "encoding/binary"

type ByteArray struct {
	data []byte
	loc int
}

func (ba *ByteArray) AppendByteArray(data []byte){
	ba.data = append(ba.data,data...)
	ba.loc+=len(data)
}

func (ba *ByteArray) AppendInt64(i64 int64) {
	long := make([]byte, 8)
	binary.LittleEndian.PutUint64(long,uint64(i64))
	ba.data = append(ba.data,long...)
	ba.loc+=4
}

func (ba *ByteArray) AppendUInt16(i16 uint16) {
	long := make([]byte, 2)
	binary.LittleEndian.PutUint16(long,uint16(i16))
	ba.data = append(ba.data,long...)
	ba.loc+=2
}

func (ba *ByteArray) GetInt64() int64 {
	result := int64(binary.LittleEndian.Uint64(ba.data[ba.loc:ba.loc+8]))
	ba.loc+=8;
	return result
}

func (ba *ByteArray) GetUInt16() uint16 {
	result := binary.LittleEndian.Uint16(ba.data[ba.loc:ba.loc+2])
	ba.loc+=2;
	return result
}

func NewByteArray(data []byte) *ByteArray {
	ba := ByteArray{}
	ba.data = data
	return &ba
}

func (ba*ByteArray) GetData()[]byte {
	return ba.data
}