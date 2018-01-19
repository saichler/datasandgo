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
	sw.StartNetworkNode(false)

	nd := net.NetNode{}
	nd.StartNetworkNode(false)

	time.Sleep(time.Second*30)

}