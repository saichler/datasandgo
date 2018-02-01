package main

import "./net"

func main() {
	node := net.Node{}
	fh := net.StringFrameHandler{}
	node.StartNetworkNode(true, fh)
}
