package main

import (
	"../net"
	"time"
)

func main(){
	sw := net.NetNode{}
	sw.StartNetworkNode(false)

	nd := net.NetNode{}
	nd.StartNetworkNode(false)

	time.Sleep(time.Second*30)
}