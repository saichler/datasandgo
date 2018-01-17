package net
import (
	"math/rand"
	"time"
	"math"
	"net"
	"log"
	"encoding/binary"
	"strconv"
	"strings"
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
	var ip int32
	ip = int32(nid.uuidLessSignificant << 32)
	return nid.uuidLessSignificant)
}

func (nid *NID) encode() []byte {
	ba := ByteArray{}
	ba.AppendInt64(nid.uuidMostSignificant)
	ba.AppendInt64(nid.uuidLessSignificant)
	ba.AppendUInt16(nid.networkId)
	ba.AppendUInt16(nid.serviceId)
	return ba.data
}

func NewNID(port int) *NID{
	newNID := NID{}
	rand.Seed(time.Now().Unix())
	newNID.uuidMostSignificant = rand.Int63n(math.MaxInt64)
	var ip int32
	ip = getIpAddress()
	newNID.uuidLessSignificant = int64(getIpAddress() << 32 + uint32(port))
	return &newNID
}

func getIpAddress() int32 {
	ifaces, err := net.Interfaces()
	if err!=nil {
		log.Fatal("Unable to access interfaces\n", err)
	}
	for _, _interface := range ifaces {
		intName := strings.ToLower(_interface.Name)
		if !strings.Contains(intName,"localhost") &&
			!strings.Contains(intName, "br") &&
				!strings.Contains(intName, "vir") {
			intAddresses, err := _interface.Addrs()
			if err!=nil {
				log.Fatal("Unable to access interface address\n", err)
			}

			for _, address := range intAddresses {
				ipaddr := address.String()
				var ipint int32
				arr := strings.Split(ipaddr,".")
				ipint = 0
				ipint += int32()
			}
		}
		return binary.BigEndian.in(ipaddr)
	}
	return 0
}