package net

import (
	"encoding/binary"
)

type ByteArray struct {
	data []byte
	loc int
}

func (ba *ByteArray) Add(data []byte){
	ba.data = append(ba.data,data...)
	ba.loc+=len(data)
}

func (ba *ByteArray) AddByteArray(data []byte){
	ba.AddUInt32(uint32(len(data)))
	ba.data = append(ba.data,data...)
	ba.loc+=len(data)+4
}

func (ba *ByteArray) GetByteArray() []byte {
	size := int(ba.GetUInt32())
	result := ba.data[ba.loc:ba.loc+size]
	ba.loc+=size
	return result
}

func (ba *ByteArray) AddInt64(i64 int64) {
	long := make([]byte, 8)
	binary.LittleEndian.PutUint64(long,uint64(i64))
	ba.data = append(ba.data,long...)
	ba.loc+=8
}

func (ba *ByteArray) GetInt64() int64 {
	result := int64(binary.LittleEndian.Uint64(ba.data[ba.loc:ba.loc+8]))
	ba.loc+=8;
	return result
}

func (ba *ByteArray) AddInt32(i32 int32) {
	long := make([]byte, 4)
	binary.LittleEndian.PutUint32(long,uint32(i32))
	ba.data = append(ba.data,long...)
	ba.loc+=4
}

func (ba *ByteArray) GetInt32() int32 {
	result := int32(binary.LittleEndian.Uint32(ba.data[ba.loc:ba.loc+4]))
	ba.loc+=4;
	return result
}

func (ba *ByteArray) AddUInt16(i16 uint16) {
	long := make([]byte, 2)
	binary.LittleEndian.PutUint16(long, i16)
	ba.data = append(ba.data,long...)
	ba.loc+=2
}

func (ba *ByteArray) GetUInt16() uint16 {
	result := binary.LittleEndian.Uint16(ba.data[ba.loc:ba.loc+2])
	ba.loc+=2;
	return result
}

func (ba *ByteArray) AddUInt32(i32 uint32) {
	long := make([]byte, 4)
	binary.LittleEndian.PutUint32(long, i32)
	ba.data = append(ba.data,long...)
	ba.loc+=4
}

func (ba *ByteArray) GetUInt32() uint32 {
	result := binary.LittleEndian.Uint32(ba.data[ba.loc:ba.loc+4])
	ba.loc+=4;
	return result
}

func (ba *ByteArray) AddBool(b bool){
	stat := make([]byte,1)
	if b {
		stat[0] = 1
	}
	ba.data = append(ba.data, stat...)
	ba.loc+=1
}

func (ba *ByteArray) GetBool() bool {
	var result bool = false
	if ba.data[ba.loc] == 1 {
		result = true
	}
	ba.loc+=1
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