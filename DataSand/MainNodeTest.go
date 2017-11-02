package main

import (
	"./network"
	"time"
)

func main() {
	sfh := network.StringFrameHandler{}
	node := network.NetNode{}
	node.FrameHandler = sfh
	node.StartNetworkNode(false)


	sfh.SendString("Hello World", node, nil)
	sfh.SendString("Hello World Again", node, nil)
	time.Sleep(time.Second*10)
	sfh.SendString("Hello World Again and Again", node, nil)
	time.Sleep(time.Second*2)
}