package net
import "../securityutil"

type Frame struct {
	source *NID
	dest *NID
	origSource *NID
	data []byte
}

type FrameHandler interface {
	HandleFrame(nNode *NetNode, frame Frame)
}

func (frame *Frame) decode (packet *Packet){
	frame.source = packet.source
	frame.dest = packet.dest
	frame.data = packet.data

	key := securityutil.SecurityKey{}
	decData, err := key.Dec(packet.data)
	if err == nil {
		frame.data = decData
	}
}

func (frame *Frame) encode() []byte {

	//@TODO add code here to split the frame into packets if it is bigger than some size
	//@TODO for now, one frame is one packet

	packet := Packet{}
	packet.source = frame.source
	packet.dest = frame.dest
	packet.data = frame.data

	key := securityutil.SecurityKey{}
	data, err := key.Enc(packet.data)
	if err == nil {
		packet.data = data
	}

	return packet.encode()
}

