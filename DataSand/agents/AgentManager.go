package agents

import (
	"../network"
	"log"
)

type AgentManager struct {
	networkNode *network.NetNode
	incoming []network.Frame
	agents  map[network.NID]Agent
}

func (am *AgentManager) Start(){
	am.networkNode = &network.NetNode{}
	am.networkNode.FrameHandler = am
	am.networkNode.StartNetworkNode(false)
	am.incoming = make([]network.Frame, 0)
	log.Println("Started Agent Manager on "+am.networkNode.Nid.String())
}

func (am *AgentManager) HandleFrame(nNode *network.NetNode, frame network.Frame){
	am.incoming = append(am.incoming, frame)
}

func (am *AgentManager)poll(){
	queue := make (chan []network.Frame)
	select {
		case queue <- am.incoming:
			f := am.incoming[0]
			am.incoming = am.incoming[:1]
			am.agents[*f.Source].HandleFrame(am.networkNode, f)
	}
}


