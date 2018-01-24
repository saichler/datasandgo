package main

import (
	"../net"
	"fmt"
	"strconv"
	"time"
)

func main(){

	ba := &net.ByteArray{}

	ba.AddInt64(876493993495)
	data := ba.GetData()
	ba2 := net.NewByteArray(data)
	fmt.Println(strconv.Itoa(int(ba2.GetInt64())))


	ip := "192.168.65.72"
	ipint32 := net.GetIpAsInt32(ip)
	ip2 := net.GetIpAsString(ipint32)
	fmt.Println(ip2)

	nid := net.NewNID(4330)
	fmt.Println(nid.String())
	data = nid.Encode()
	ba = net.NewByteArray(data)

	nid2 := nid.Decode(ba)
	fmt.Println(nid2)

	packet := net.Packet{}
	packet.SetSource(nid2)
	data = packet.Encode()
	ba = net.NewByteArray(data)
	packet2 := packet.HeaderDecode(ba)
	packet.DataDecode(ba)
	fmt.Println(packet2.String())

	sw := net.Node{}
	swfh := net.StringFrameHandler{}
	sw.StartNetworkNode(false, swfh)

	nd1 := net.Node{}
	ndfh1 := net.StringFrameHandler{}
	nd1.StartNetworkNode(false, ndfh1)

	time.Sleep(time.Second*2)
/*
	fmt.Println("nid:"+nd1.GetNID().String())
	ndfh1.SendString("Hello World",&nd1,nd1.GetSwitchNID())
	time.Sleep(time.Second*2)
	longString := ""
	for i:=0;i<net.MAX_PACKET_SIZE+100;i++ {
		longString+="A"
	}

	ndfh1.SendString(longString,&nd1,nd1.GetSwitchNID())

	time.Sleep(time.Second*2)

	longString = ""
	for i:=0;i<net.MAX_PACKET_SIZE*300+7;i++ {
		longString+="A"
	}

	ndfh1.SendString(longString,&nd1,nd1.GetSwitchNID())

	time.Sleep(time.Second*2)
*/
	nd2 := net.Node{}
	ndfh2 := net.StringFrameHandler{}
	nd2.StartNetworkNode(false, ndfh2)

	time.Sleep(time.Second*2)

	ndfh2.SendString("Hello Adjacent",&nd2,nd1.GetNID())

	time.Sleep(time.Second*2)

	longString := ""
	for i:=0;i<net.MAX_PACKET_SIZE*300+7;i++ {
		longString+="B"
	}

	ndfh2.SendString(longString,&nd2,nd1.GetNID())

	time.Sleep(time.Second*60)
}