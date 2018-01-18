package net

type Packet struct {
	source *NID
	dest *NID
	origSource *NID
	packetId uint32
	part uint32
	multipart bool
	priority uint8
	data []byte
}

func (p *Packet) encode() []byte {
	ba := ByteArray{}
	ba.AppendByteArray(p.source.encode())
	ba.AppendByteArray(p.dest.encode())
	ba.AppendByteArray(p.origSource.encode())
	ba.AppendUInt32(p.packetId)
	ba.AppendUInt32(p.part)
	
	result := make([]byte,12+12+12+4+4+2+2+len(p.data))
	nidSourceData := p.source.encode()
	copy(result[0:12],nidSourceData[:])
	return result
}

func (p *Packet) decode(data []byte) {
}
