package agents

import "../network"

type Agent interface {
	HandleFrame(nNode *network.NetNode, frame network.Frame)
}
