package net

type Packet struct {
	source *NID
	dest *NID
	origSource *NID
	frameID uint32
	pnum uint32
	multipart bool
	priority uint16
	data []byte
}

var nidDecoder = NID{}

const (
	MAX_PACKET_SIZE = 512
)

func (p *Packet) Encode() []byte {
	ba := ByteArray{}
	ba.Add(p.source.Encode())
	if p.dest == nil {
		p.dest = NewNID(0)
	}
	ba.Add(p.dest.Encode())
	if p.origSource == nil {
		p.origSource = NewNID(0)
	}
	ba.Add(p.origSource.Encode())
	ba.AddUInt32(p.frameID)
	ba.AddUInt32(p.pnum)
	ba.AddBool(p.multipart)
	ba.AddUInt16(p.priority)
	ba.AddByteArray(p.data)

	return ba.data
}

func (p *Packet)HeaderDecode(ba *ByteArray) *Packet{
	packet := Packet{}
	packet.source = nidDecoder.Decode(ba)
	packet.dest = nidDecoder.Decode(ba)
	return &packet
}

func (packet *Packet)DataDecode(ba *ByteArray){
	packet.origSource = nidDecoder.Decode(ba)
	packet.frameID = ba.GetUInt32()
	packet.pnum = ba.GetUInt32()
	packet.multipart = ba.GetBool()
	packet.priority = ba.GetUInt16()
	packet.data = ba.GetByteArray()
}

func (p *Packet) String() string {
	return p.source.String()
}

func (p *Packet) SetSource(nid *NID){
	p.source = nid
}
