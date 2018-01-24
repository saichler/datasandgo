package net

import (
	"net"
	"log"
	"encoding/binary"
	"strconv"
	"strings"
)

const (
	SWITCH_PORT = 52000
	MAX_PORT = 54000
)

type Node struct {
	nid *NID
	links map[NID]net.Conn
	frameHandler FrameHandler
	isSwitch bool
}

var packetDecoder = Packet{}

func (node *Node) StartNetworkNode(service bool, handler FrameHandler){
	node.links = make(map[NID]net.Conn)
	node.frameHandler = handler
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
		if port != SWITCH_PORT {
			port--
		}
		node.nid = NewNID(port)
		log.Println("Bounded to port "+node.nid.String())
		node.isSwitch = isSwitch
		if !isSwitch {
			node.uplinkToSwitch()
		}
	}
	if service {
		node.waitForlinks(socket)
	} else {
		go node.waitForlinks(socket)
	}
}

func (node *Node)waitForlinks(socket net.Listener){
	//infinit loop to accept connections
	for {
		connection, error := socket.Accept()
		if error != nil {
			log.Fatal("Failed to accept a new connection from socket: ", error)
			return
		}
		//start a new connection
		go node.newConnection(connection)
	}
	log.Fatal("Server Socket was closed!")
}

func (node *Node)newConnection(c net.Conn){
	log.Println("Connected to: "+c.RemoteAddr().String())

	node.handshake(c)

	chanSize := make(chan []byte)
	chanError := make(chan error)

	for {
		data := node.singlePacketRead(c, chanSize, chanError)
		if data != nil {
			node.handlePacket(data)
		} else {
			break;
		}
	}
}

func (node *Node)singlePacketRead(c net.Conn, chanSize chan []byte, chanError chan error)([]byte){
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
		node.unregisterLink(c)
		c.Close()
		log.Println("Connection of "+c.RemoteAddr().String()+" was closed!", err)
	}
	return nil
}

func (node *Node)unregisterLink(c net.Conn){
	var keyToRemove NID
	for key, value := range node.links {
		if(value == c){
			keyToRemove = key
			break;
		}
	}
	node.links[keyToRemove]=nil
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

func (node *Node)handlePacket(data []byte){
	ba := NewByteArray(data)
	packet := packetDecoder.HeaderDecode(ba)
	if packet.dest.equal(node.nid) {
		packet.DataDecode(ba)
		frame := Frame{}
		frame.Decode(packet)
		if frame.complete {
			node.frameHandler.HandleFrame(node, &frame)
		}
	} else {
		node.sendBytes(packet, data)
	}
}

func (node *Node)handshake(c net.Conn){
	log.Println("Handshake")
	packet := Packet{}
	packet.source = node.nid
	data := packet.Encode()
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	c.Write(size)
	c.Write(data)

	chanSize := make(chan []byte)
	chanError := make(chan error)

	data = node.singlePacketRead(c, chanSize, chanError)
	ba := NewByteArray(data)
	p := packetDecoder.HeaderDecode(ba)

	log.Println("handshaked with nid:"+p.source.String())

	node.links[*p.source] = c
}

func (node *Node)uplinkToSwitch() {
	switchPortString := strconv.Itoa(SWITCH_PORT)
	c, e := net.Dial("tcp", "127.0.0.1:"+switchPortString)
	if e != nil {
		log.Fatal("Failed to open connection to switch: ", e)
	}

	go node.newConnection(c)
}

func (node *Node)sendBytes(packet *Packet, data []byte){
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	var c net.Conn
	if !node.isSwitch {
		swNID := node.GetSwitchNID()
		c = node.links[*swNID]
	} else {
		c = node.links[*packet.dest]
	}
	log.Println("Sending from "+node.nid.String() +" to "+packet.dest.String())
	if c==nil{
		for key,_ := range node.links {
			log.Println("NID1:"+key.String()+"\nNID2:"+packet.dest.String())
		}
		log.Fatal("Invalid Connection to :"+packet.dest.String())
	}
	c.Write(size)
	c.Write(data)
}

func (node *Node)send(packet *Packet){
	data := packet.Encode()
	node.sendBytes(packet, data)
}

func (node *Node)Send(frame *Frame) {
	packets := frame.Encode()
	for i:=0;i<len(packets); i++ {
		node.send(packets[i])
	}
}

func (node *Node) GetSwitchNID() *NID {
	for key, _ := range node.links {
		if strings.Contains(key.String(),"52000") {
			return &key
		}
	}
	return nil
}

func (node *Node) GetNID () *NID {
	return node.nid
}