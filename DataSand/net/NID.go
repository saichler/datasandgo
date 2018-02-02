package net
import (
	"math/rand"
	"time"
	"math"
	"net"
	"log"
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
	ip := int32(nid.uuidLessSignificant >> 32)
	port := int(nid.uuidLessSignificant - ((nid.uuidLessSignificant >> 32) << 32))
	return strconv.Itoa(int(nid.uuidMostSignificant))+":"+GetIpAsString(ip)+":"+strconv.Itoa(port)
}

func (nid *NID) Encode() []byte {
	ba := ByteArray{}
	ba.AddInt64(nid.uuidMostSignificant)
	ba.AddInt64(nid.uuidLessSignificant)
	ba.AddUInt16(nid.networkId)
	ba.AddUInt16(nid.serviceId)
	return ba.data
}

func (n *NID)Decode(ba *ByteArray) *NID{
	nid := NewNID(0)
	nid.uuidMostSignificant = ba.GetInt64()
	nid.uuidLessSignificant = ba.GetInt64()
	nid.networkId = ba.GetUInt16()
	nid.serviceId = ba.GetUInt16()
	return nid
}

func NewNID(port int) *NID{
	newNID := NID{}
	rand.Seed(time.Now().Unix())
	newNID.uuidMostSignificant = rand.Int63n(math.MaxInt64)
	var ip int32
	ip = getIpAddress()
	newNID.uuidLessSignificant = int64(ip) << 32 + int64(port)
	return &newNID
}

func getIpAddress() int32 {
	ifaces, err := net.Interfaces()
	if err!=nil {
		log.Fatal("Unable to access interfaces\n", err)
	}
	for _, _interface := range ifaces {
		intName := strings.ToLower(_interface.Name)
		if !strings.Contains(intName,"lo") &&
			!strings.Contains(intName, "br") &&
				!strings.Contains(intName, "vir") {
			intAddresses, err := _interface.Addrs()
			if err!=nil {
				log.Fatal("Unable to access interface address\n", err)
			}

			for _, address := range intAddresses {
				ipaddr := address.String()
				return GetIpAsInt32(ipaddr)
			}
		}
	}
	return 0
}

func GetIpAsString( ip int32) string {
	a := strconv.FormatInt(int64((ip>>24)&0xff), 10)
	b := strconv.FormatInt(int64((ip>>16)&0xff), 10)
	c := strconv.FormatInt(int64((ip>>8)&0xff), 10)
	d := strconv.FormatInt(int64(ip & 0xff), 10)
	return a + "." + b + "." + c + "." + d
}

func GetIpAsInt32(ipaddr string) int32 {
	var ipint int32
	arr := strings.Split(ipaddr,".")
	ipint = 0
	a,_ := strconv.Atoi(arr[0])
	b,_ := strconv.Atoi(arr[1])
	c,_ := strconv.Atoi(arr[2])
	d,_ := strconv.Atoi(strings.Split(arr[3],"/")[0])
	ipint += int32(a) << 24
	ipint += int32(b) << 16
	ipint += int32(c) << 8
	ipint += int32(d)
	return ipint
}

func FromString(str string) *NID {
	nid := NID{}
	index := strings.Index(str,":")
	mostString :=  str[0:index]
	lessString := str[index+1:len(str)]
	index1 := strings.Index(lessString,":")

	mostUUID,_ := strconv.Atoi(mostString)
	nid.uuidMostSignificant = int64(mostUUID)

	ip := GetIpAsInt32(lessString[0:index1])
	port,_ := strconv.Atoi(lessString[index1+1:len(lessString)])

	nid.uuidLessSignificant = int64(ip) << 32 + int64(port)
	return &nid
}

func (nid *NID) equal (other *NID) bool {
	return  nid.uuidMostSignificant == other.uuidMostSignificant &&
			nid.uuidLessSignificant == other.uuidLessSignificant &&
			nid.networkId == other.networkId &&
			nid.serviceId == other.serviceId
}

func (nid *NID) sameMachine(other *NID) bool {
	myip := int32(nid.uuidLessSignificant >> 32)
	otherip := int32(other.uuidLessSignificant >> 32)
	return myip == otherip
}

func (nid *NID) getHostID() int32 {
	return int32(nid.uuidLessSignificant >> 32)
}