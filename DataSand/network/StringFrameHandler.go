package network

import "log"

type StringFrameHandler struct {
}

const (
	REQUEST = 1
	REPLY = 2;
)

func (sfh StringFrameHandler) HandleFrame(nNode *NetNode, frame Frame){
	if frame.FrameType == REQUEST {
		log.Println("Request: "+string(frame.Data))
		sfh.ReplyString(string(frame.Data), nNode, frame.Source)
	} else {
		log.Println("Reply: "+string(frame.Data))
	}
}

func (sfh StringFrameHandler)SendString(str string, nNode *NetNode, dest *NID){
	log.Println("Sending "+str)
	frame := Frame{}
	if dest==nil {
		frame.Dest = &nNode.switchNid
	} else {
		frame.Dest = dest
	}
	nid := NID{}
	nid.Uuid = nNode.Nid.Uuid
	frame.Source = &nid
	frame.Data = []byte(str)
	frame.FrameType = REQUEST;
	nNode.Send(frame)
}

func (sfh StringFrameHandler)ReplyString(str string, nNode *NetNode, dest *NID){
	frame := Frame{}
	if dest==nil {
		frame.Dest = &nNode.switchNid
	} else {
		frame.Dest = dest
	}
	nid := NID{}
	nid.Uuid = nNode.Nid.Uuid
	frame.Source = &nid
	frame.Data = []byte(str)
	frame.FrameType = REPLY;
	nNode.Send(frame)
}