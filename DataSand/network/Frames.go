package network

type FrameHandler interface {
	HandleFrame(nNode NetNode, frame Frame)
}
