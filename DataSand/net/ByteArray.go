package net

import "encoding/binary"

type ByteArray struct {
	data []byte
	loc int
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