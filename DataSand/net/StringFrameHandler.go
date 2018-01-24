package net

import "log"

type StringFrameHandler struct {
}

type Protocol struct {
	op uint32
	data string
}

const (
	REQUEST = 1
	REPLY = 2;
)

func getData(frame *Frame) *Protocol {
	ba := ByteArray{}
	ba.data = frame.data
	protocol := Protocol{}
	protocol.op = ba.GetUInt32()
	protocol.data = ba.GetString()
	return &protocol
}

func (sfh StringFrameHandler) HandleFrame(node *Node, frame *Frame){
	protocol := getData(frame)
	if protocol.op == REQUEST {
		log.Println("Request: "+protocol.data)
		sfh.ReplyString(protocol.data, node, frame.source)
	} else {
		log.Println("Reply: "+protocol.data)
	}
}

func (sfh StringFrameHandler)SendString(str string, node *Node, dest *NID){
	frame := NewFrame()
	if dest==nil {
		frame.dest = node.GetSwitchNID()
	} else {
		frame.dest = dest
	}

	frame.source = node.nid
	ba := ByteArray{}
	ba.AddUInt32(REQUEST)
	ba.AddString(str)
	frame.data = ba.data

	node.Send(frame)
}

func (sfh StringFrameHandler)ReplyString(str string, node *Node, dest *NID){
	frame := NewFrame()
	if dest==nil {
		frame.dest = node.GetSwitchNID()
	} else {
		frame.dest = dest
	}

	frame.source = node.nid
	ba := ByteArray{}
	ba.AddUInt32(REPLY)
	ba.AddString(str)
	frame.data = ba.data

	node.Send(frame)
}