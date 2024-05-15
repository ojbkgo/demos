package protocol

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const HeaderSize = 12
const MTU = 1500

type UdpPackage struct {
	Header      [HeaderSize]byte // 4 bytes for length, 4 bytes for src ip, 4 bytes for dest ip
	Payload     []byte
	SrcIp       string
	DestIp      string
	PayloadSize uint
}

func NewUdpPackage(src, dest string, payload []byte) *UdpPackage {
	p := &UdpPackage{
		Payload: make([]byte, len(payload)),
	}

	copy(p.Payload, payload)

	p.SrcIp = src
	p.DestIp = dest

	// set header
	length := len(payload)
	p.Header[0] = byte(length)
	p.Header[1] = byte(length >> 8)
	p.Header[2] = byte(length >> 16)
	p.Header[3] = byte(length >> 24)

	srcIpByte := p.convertIPInteger(src)
	destIpByte := p.convertIPInteger(dest)
	copy(p.Header[4:8], srcIpByte)
	copy(p.Header[8:12], destIpByte)

	return p
}

func (p *UdpPackage) convertIPInteger(ip string) []byte {
	ipParts := strings.Split(ip, ".")
	res := make([]byte, 4)
	for i := 0; i < 4; i++ {
		ii, _ := strconv.Atoi(ipParts[i])
		res[i] = byte(ii)
	}
	return res
}

func (p *UdpPackage) convertIPString(ip []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (p *UdpPackage) Serialize() []byte {
	data := make([]byte, len(p.Header)+len(p.Payload))
	copy(data, p.Header[:])
	copy(data[HeaderSize:], p.Payload)
	return data
}

func (p *UdpPackage) ReadPackage(conn *net.UDPConn) (*net.UDPAddr, error) {
	udpPayload := make([]byte, MTU)

	n, addr, err := conn.ReadFromUDP(udpPayload)
	if err != nil {
		return addr, err
	}
	header := udpPayload[:HeaderSize]
	// header byte 转 16进制
	fmt.Println("header byte to hex:")
	for i := 0; i < len(header); i++ {
		fmt.Printf("%x ", header[i])
	}

	fmt.Println("read header size:", n, addr.String())

	p.PayloadSize = uint(header[0]) | uint(header[1])<<8 | uint(header[2])<<16 | uint(header[3])<<24
	p.SrcIp = p.convertIPString(header[4:8])
	p.DestIp = p.convertIPString(header[8:12])
	fmt.Println("parse header success....: ", p.PayloadSize, p.SrcIp, p.DestIp)

	p.Payload = make([]byte, p.PayloadSize)
	copy(p.Payload, udpPayload[HeaderSize:HeaderSize+int(p.PayloadSize)])

	return addr, nil
}
