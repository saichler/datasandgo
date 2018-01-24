package net

type Frame struct {
	frameID uint32
	source *NID
	dest *NID
	origSource *NID
	data []byte
	complete bool
}

type MultiPart struct {
	frameID uint32
	packets map[uint32]*Packet
	totalExpectedPackets uint32
}

var pending map[uint32]*MultiPart = make(map[uint32]*MultiPart)

type FrameHandler interface {
	HandleFrame(nNode *NetNode, frame *Frame)
}

func (frame *Frame) Decode (packet *Packet){
	frame.source = packet.source
	frame.dest = packet.dest

	if packet.multipart {
		var multiPart *MultiPart
		multiPart = pending[packet.frameID]
		if multiPart == nil {
			multiPart = &MultiPart{}
			multiPart.packets = make(map[uint32]*Packet)
			pending[packet.frameID] = multiPart
		}
		multiPart.packets[packet.pnum] = packet
		if multiPart.totalExpectedPackets == 0 && packet.pnum == 0 {
			ba := ByteArray{}
			ba.data = packet.data
			multiPart.totalExpectedPackets = ba.GetUInt32()
		}
		if multiPart.totalExpectedPackets>0 && len(multiPart.packets) == int(multiPart.totalExpectedPackets) {
			frameData := make([]byte,0)
			for i:=1;i<int(multiPart.totalExpectedPackets);i++ {
				frameData = append(frameData, multiPart.packets[uint32(i)].data...)
			}
			/* decrypt here
			key := securityutil.SecurityKey{}
			decData, err := key.Dec(packet.data)
			if err == nil {
				frame.data = decData
			}*/
			frame.data = frameData
			pending[packet.frameID] = nil
			frame.complete = true;
		}
	} else {
		/* decrypt here
		key := securityutil.SecurityKey{}
		decData, err := key.Dec(packet.data)
		if err == nil {
			frame.data = decData
		}*/
		frame.data = packet.data
		frame.complete = true
	}
}

func (frame *Frame) Encode() []*Packet {

	frameData := frame.data

	/* encrypt here
key := securityutil.SecurityKey{}
data, err := key.Enc(packet.data)
if err == nil {
	packet.data = data
}*/

	if len(frameData)> MAX_PACKET_SIZE {
		totalParts := len(frameData)/MAX_PACKET_SIZE
		left := len(frame.data) - totalParts*MAX_PACKET_SIZE
		if left>0 {
			totalParts++
		}
		totalParts++

		result := make([]*Packet,totalParts)

		ba := ByteArray{}
		ba.AddUInt32(uint32(totalParts))

		packet := Packet{}
		packet.source = frame.source
		packet.dest = frame.dest
		packet.data = ba.data
		packet.multipart = true
		packet.pnum = 0
		packet.frameID = frame.frameID
		result[0] = &packet

		for i:=0;i<totalParts-1;i++ {
			loc := i*MAX_PACKET_SIZE
			packet := Packet{}
			packet.source = frame.source
			packet.dest = frame.dest
			packet.frameID = frame.frameID
			if i<totalParts-2 || left==0 {
				packet.data = frameData[loc:loc+MAX_PACKET_SIZE]
			} else {
				packet.data = frameData[loc:loc+left]
			}

			packet.multipart = true
			packet.pnum = uint32(i+1)
			result[i+1] = &packet
		}
		return result
	} else {
		result := make([]*Packet,1)
		packet := Packet{}
		packet.source = frame.source
		packet.dest = frame.dest
		packet.data = frame.data
		packet.frameID = frame.frameID
		packet.multipart = false
		result[0] = &packet
		return result
	}
}