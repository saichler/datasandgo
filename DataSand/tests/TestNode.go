package main

import (
	"../net"
	"fmt"
	"strconv"
)

func main(){
	ba := net.ByteArray{}

	ba.AppendInt64(876493993495)
	data := ba.GetData()
	ba2 := net.NewByteArray(data)
	fmt.Println(strconv.Itoa(int(ba2.GetInt64())))
	/*
	sw := net.NetNode{}
	sw.StartNetworkNode(false)

	nd := net.NetNode{}
	nd.StartNetworkNode(false)

	time.Sleep(time.Second*30)
	*/
}