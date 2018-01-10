package net
import (
	"math/rand"
	"time"
	"math"
	"net"
	"log"
	"encoding/binary"
	"fmt"
	"strconv"
)

type NID struct {
	uuidMostSignificant int64
	uuidLessSignificant int64
	networkId uint16
	serviceId uint16
}

func (nid *NID) getMostSignificant() int64 {
	return nid.uuidMostSignificant
}

func (nid *NID) getLessSignificant() int64 {
	return nid.uuidLessSignificant
}

func (nid *NID) getNetworkId() uint16 {
	return nid.networkId
}

func (nid *NID) getServiceId() uint16 {
	return nid.serviceId
}

func (nid *NID) String() string {
	return string(nid.uuidLessSignificant)
}

func (nid *NID) encode() []byte {
	ba := ByteArray{}
	ba.AppendInt64(nid.uuidMostSignificant)
	ba.AppendInt64(nid.uuidLessSignificant)
	ba.AppendUInt16(nid.networkId)
	ba.AppendUInt16(nid.serviceId)
	fmt.Println("size="+strconv.Itoa(len(ba.data)))
	return ba.data
}

func NewNID(port int) *NID{
	newNID := NID{}
	rand.Seed(time.Now().Unix())
	newNID.uuidMostSignificant = rand.Int63n(math.MaxInt64)
	newNID.uuidLessSignificant = int64(getIpAddress()+uint32(port))
	return &newNID
}

func getIpAddress() uint32 {
	ifaces, err := net.Interfaces()
	if err!=nil {
		log.Fatal("Unable to access interfaces\n", err)
	}
	for _, _interface := range ifaces {
		intAddresses, err := _interface.Addrs()
		if err!=nil {
			log.Fatal("Unable to access interface address\n", err)
		}

		ipaddr := net.IP{}

		for _, address := range intAddresses {
			switch value := address.(type) {
			case *net.IPNet:
				ipaddr = value.IP
			case *net.IPAddr:
				ipaddr = value.IP
			}
		}
		return binary.BigEndian.Uint32(ipaddr)
	}
	return 0
}