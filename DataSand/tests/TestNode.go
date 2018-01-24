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
	packet2 := packet.Decode(ba)
	fmt.Println(packet2.String())

	sw := net.NetNode{}
	swfh := net.StringFrameHandler{}
	sw.StartNetworkNode(false, swfh)

	nd := net.NetNode{}
	ndfh := net.StringFrameHandler{}
	nd.StartNetworkNode(false, ndfh)

	time.Sleep(time.Second*2)

	fmt.Println("nid:"+nd.GetNID())
	ndfh.SendString("Hello World",&nd,nd.GetSwitchNID())
	time.Sleep(time.Second*2)
	longString := ""
	for i:=0;i<net.MAX_PACKET_SIZE+100;i++ {
		longString+="A"
	}

	ndfh.SendString(longString,&nd,nd.GetSwitchNID())

	time.Sleep(time.Second*60)
}