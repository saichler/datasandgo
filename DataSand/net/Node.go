package net

import (
	"net"
	"log"
	"encoding/binary"
	"strconv"
)

const (
	SWITCH_PORT = 52000
	MAX_PORT = 54000
)

type NetNode struct {
	nid *NID
	links map[*NID]net.Conn
	frameHandler FrameHandler
	isSwitch bool
}

func (nNode *NetNode) StartNetworkNode(service bool){
	nNode.links = make(map[*NID]net.Conn)
	var port = SWITCH_PORT
	var portString = strconv.Itoa(port)
	isSwitch := true
	log.Println("Trying to bind to switch port "+portString+".");
	socket, error := net.Listen("tcp", ":"+portString)
	//
	if error!=nil {
		isSwitch = false
		for ; port < MAX_PORT && error != nil; port++ {
			portString = strconv.Itoa(port)
			log.Println("Trying to bind to port "+portString+".")
			s, e := net.Listen("tcp", ":"+portString)
			error = e
			socket = s
		}
		log.Println("Successfuly binded to port "+portString)
	}

	if error != nil {
		log.Fatal("Unable to bind to any of the ports.: ", error)
		return
	} else {
		nNode.nid = NewNID(port)
		log.Println("Bounded to port "+portString)
		nNode.isSwitch = isSwitch
		if !isSwitch {
			nNode.uplinkToSwitch()
		}
	}
	if service {
		nNode.waitForlinks(socket)
	} else {
		go nNode.waitForlinks(socket)
	}
}

func (nNode *NetNode)waitForlinks(socket net.Listener){
	//infinit loop to accept connections
	for {
		connection, error := socket.Accept()
		if error != nil {
			log.Fatal("Failed to accept a new connection from socket: ", error)
			return
		}
		//start a new connection
		go nNode.newConnection(connection)
	}
	log.Fatal("Server Socket was closed!")
}

func (nNode *NetNode)newConnection(c net.Conn){
	log.Println("Connected to: "+c.RemoteAddr().String())

	nNode.handshake(c)

	chanSize := make(chan []byte)
	chanError := make(chan error)

	for {
		data := nNode.singlePacketRead(c, chanSize, chanError)
		if data != nil {
			nNode.handlePacket(data)
		} else {
			break;
		}
	}
}

func (nNode *NetNode)singlePacketRead(c net.Conn, chanSize chan []byte, chanError chan error)([]byte){
	go readDataSize(c, chanSize, chanError)

	select {
	case sizeData := <-chanSize:
		dataSize:= int(binary.LittleEndian.Uint32(sizeData))
		data, err :=readData(c, dataSize)
		if data!=nil {
			return data
		} else if err!=nil {
			log.Fatal("Failed to read data from "+c.RemoteAddr().String()+": ",err)
			break;
		}
	case err := <-chanError:
		nNode.unregisterLink(c)
		c.Close()
		log.Println("Connection of "+c.RemoteAddr().String()+" was closed!", err)
	}
	return nil
}

func (nNode *NetNode)unregisterLink(c net.Conn){
	keyToRemove := &NID{}
	for key, value := range nNode.links {
		if(value == c){
			keyToRemove = key
			break;
		}
	}
	nNode.links[keyToRemove]=nil
}

func readDataSize(c net.Conn, chanSize chan []byte, chanError chan error){
	dataSizeInBytes := make([]byte, 4)
	_,e := c.Read(dataSizeInBytes)

	if e != nil {
		log.Println("Failed to read data size, closing connection!", e)
		chanError<-e
	} else {
		chanSize<-dataSizeInBytes
	}
}

func readData(c net.Conn, size int) ([]byte, error) {
	data := make([]byte, size)
	_,e := c.Read(data)
	if e != nil {
		log.Fatal("Failed to read data ", e)
		return nil, e
	}
	return data, nil
}

func (nNode *NetNode)handlePacket(data []byte){
	log.Println("Handle Packet")
	packet := Packet{}
	packet.decode(data)

	//@TODO add code here to collect the packets to a set
	//@TODO so when all packet have arrived for a frame
	//@TODO combine the packets data and deserialize a frame
	//@TODO for now, one packet is one frame

	frame := Frame{}
	frame.decode(&packet)

	nNode.frameHandler.HandleFrame(nNode, frame)
}

func (nNode *NetNode)handshake(c net.Conn){
	log.Println("Handshake")
	packet := Packet{}
	packet.source = nNode.nid
	data := packet.encode()
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	c.Write(size)
	c.Write(data)

	chanSize := make(chan []byte)
	chanError := make(chan error)

	data = nNode.singlePacketRead(c, chanSize, chanError)

	p := Packet{}
	p.decode(data)
	nNode.links[p.source] = c
}

func (nNode *NetNode)uplinkToSwitch() {
	switchPortString := strconv.Itoa(SWITCH_PORT)
	c, e := net.Dial("tcp", "127.0.0.1:"+switchPortString)
	if e != nil {
		log.Fatal("Failed to open connection to switch: ", e)
	}

	go nNode.newConnection(c)
}

func (nNode *NetNode)send(data[] byte, nid *NID){
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	c := nNode.links[nid]
	c.Write(size)
	c.Write(data)
}

func (nNode *NetNode)Send(frame Frame) {
	nNode.send(frame.encode(), frame.dest)
}