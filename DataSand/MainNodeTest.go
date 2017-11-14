package main

import (
	"./network"
	"time"
	"log"
)

func main() {
	sfh := network.StringFrameHandler{}
	node := network.NetNode{}
	node.FrameHandler = sfh
	node.StartNetworkNode(false)

	time.Sleep(time.Second*5)
	log.Println("Sending first packet")
	sfh.SendString("Hello World", &node, nil)
	time.Sleep(time.Second*5)
	sfh.SendString("Hello World Again", &node, nil)
	time.Sleep(time.Second*5)
	sfh.SendString("Hello World Again and Again", &node, nil)
	time.Sleep(time.Second*5)
}