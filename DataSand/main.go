package main

import "./net"

func main() {
	node := net.Node{}
	fh := net.StringFrameHandler{}
	node.StartNetworkNode(true, fh)

	node1 := net.Node{}
	fh1 := net.StringFrameHandler{}
	node1.StartNetworkNode(true, fh1)

}
