package net

import (
	"net"
	"log"
	"encoding/binary"
	"strconv"
	"strings"
	"time"
)

const (
	SWITCH_PORT = 52000
	MAX_PORT    = 54000
)

type Node struct {
	nid          *NID
	interfaces   map[NID]net.Conn
	peers map[int32]net.Conn
	peersNid map[int32]*NID
	tunnels      map[NID]*NID
	frameHandler FrameHandler
	isSwitch     bool
}

var packetDecoder = Packet{}

func (node *Node) StartNetworkNode(service bool, handler FrameHandler) {
	node.interfaces = make(map[NID]net.Conn)
	node.peers = make(map[int32]net.Conn)
	node.peersNid = make(map[int32]*NID)
	node.frameHandler = handler
	var port = SWITCH_PORT
	var portString = strconv.Itoa(port)
	isSwitch := true
	log.Println("Trying to bind to switch port " + portString + ".");
	socket, error := net.Listen("tcp", ":"+portString)
	//
	if error != nil {
		isSwitch = false
		for ; port < MAX_PORT && error != nil; port++ {
			portString = strconv.Itoa(port)
			log.Println("Trying to bind to port " + portString + ".")
			s, e := net.Listen("tcp", ":"+portString)
			error = e
			socket = s
		}
		log.Println("Successfuly binded to port " + portString)
	}

	if error != nil {
		log.Fatal("Unable to bind to any of the ports.: ", error)
		return
	} else {
		if port != SWITCH_PORT {
			port--
		}
		node.nid = NewNID(port)
		log.Println("Bounded to port " + node.nid.String())
		node.isSwitch = isSwitch
		if !isSwitch {
			node.uplinkToSwitch()
		}
	}
	if service {
		node.waitForlinks(socket)
	} else {
		go node.waitForlinks(socket)
		time.Sleep(time.Second)
	}
}

func (node *Node) waitForlinks(socket net.Listener) {
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

func (node *Node) newConnection(c net.Conn) {
	log.Println("Connected to: " + c.RemoteAddr().String())

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

func (node *Node) singlePacketRead(c net.Conn, chanSize chan []byte, chanError chan error) ([]byte) {
	go readDataSize(c, chanSize, chanError)

	select {
	case sizeData := <-chanSize:
		dataSize := int(binary.LittleEndian.Uint32(sizeData))
		data, err := readData(c, dataSize)
		if data != nil {
			return data
		} else if err != nil {
			log.Fatal("Failed to read data from "+c.RemoteAddr().String()+": ", err)
			break;
		}
	case err := <-chanError:
		node.unregisterLink(c)
		c.Close()
		log.Println("Connection of "+c.RemoteAddr().String()+" was closed!", err)
	}
	return nil
}

func (node *Node) unregisterLink(c net.Conn) {
	for key, value := range node.interfaces {
		if (value == c) {
			node.interfaces[key] = nil
			break;
		}
	}
	for key, value := range node.peers {
		if (value == c) {
			node.peers[key] = nil
			break;
		}
	}
}

func readDataSize(c net.Conn, chanSize chan []byte, chanError chan error) {
	dataSizeInBytes := make([]byte, 4)
	_, e := c.Read(dataSizeInBytes)

	if e != nil {
		log.Println("Failed to read data size, closing connection!", e)
		chanError <- e
	} else {
		chanSize <- dataSizeInBytes
	}
}

func readData(c net.Conn, size int) ([]byte, error) {
	data := make([]byte, size)
	_, e := c.Read(data)
	if e != nil {
		log.Fatal("Failed to read data ", e)
		return nil, e
	}
	return data, nil
}

func (node *Node) handlePacket(data []byte) {
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
		if node.nid.sameMachine(packet.dest) {
			node.sendBytes(packet, data)
		} else {
			log.Fatal("External Switching is not supported yet")
		}
	}
}

func (node *Node) handshake(c net.Conn) {
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

	log.Println("handshaked with nid:" + p.source.String())

	hostID := p.source.getHostID()

	if node.nid.getHostID()==hostID {
		node.interfaces[*p.source] = c
	} else {
		node.peers[hostID] = c
		node.peersNid[hostID] = p.source
	}
}

func (node *Node) uplinkToSwitch() {
	switchPortString := strconv.Itoa(SWITCH_PORT)
	c, e := net.Dial("tcp", "127.0.0.1:"+switchPortString)
	if e != nil {
		log.Fatal("Failed to open connection to switch: ", e)
	}

	go node.newConnection(c)
}

func (node *Node) Uplink(host string) {
	switchPortString := strconv.Itoa(SWITCH_PORT)
	c, e := net.Dial("tcp", host+":"+switchPortString)
	if e != nil {
		log.Fatal("Failed to open connection to host: "+host, e)
	}

	go node.newConnection(c)

	hostID := GetIpAsInt32(host)

	for node.peers[hostID]==nil {
		time.Sleep(time.Second)
	}
}

func (node *Node) sendBytes(packet *Packet, data []byte) {
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	var c net.Conn
	hostID := packet.dest.getHostID()
	myHostID := node.nid.getHostID()

	if !node.isSwitch {
		swNID := node.GetSwitchNID()
		c = node.interfaces[*swNID]
	} else if hostID == myHostID {
		c = node.interfaces[*packet.dest]
	} else {
		c = node.peers[packet.dest.getHostID()]
	}
	//log.Println("Sending from " + node.nid.String() + " to " + packet.dest.String())
	if c == nil {
		for key, _ := range node.interfaces {
			log.Println("NID1:" + key.String() + "\nNID2:" + packet.dest.String())
		}
		var nid *NID
		nid.String()
		log.Fatal("Invalid Connection to :" + packet.dest.String())
	}
	c.Write(size)
	c.Write(data)
}

func (node *Node) send(packet *Packet) {
	data := packet.Encode()
	node.sendBytes(packet, data)
}

func (node *Node) Send(frame *Frame) {
	packets := frame.Encode()
	for i := 0; i < len(packets); i++ {
		node.send(packets[i])
	}
}

func (node *Node) GetSwitchNID() *NID {
	for key, _ := range node.interfaces {
		if strings.Contains(key.String(), "52000") {
			return &key
		}
	}
	return nil
}

func (node *Node) GetNID() *NID {
	return node.nid
}

func (node *Node) GetNodeSwitch(host string) *NID {
	hostID := GetIpAsInt32(host)
	return node.peersNid[hostID]
}