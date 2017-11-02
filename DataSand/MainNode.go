package main

import (
	"./network"
)
func main() {
	sfh := network.StringFrameHandler{}
	netnode:=network.NetNode{}
	netnode.FrameHandler = sfh
	netnode.StartNetworkNode(true)
}
