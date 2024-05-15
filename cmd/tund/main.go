package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/net/ipv4"

	"github.com/ojbkgo/demos/pkg/protocol"
	"github.com/ojbkgo/demos/pkg/utils/nettool"
)

func main() {
	// 监听udp 8299
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8299")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	netTool := nettool.NewMacNetTool()
	ifce, err := netTool.CreateTun("tun99")
	if err != nil {
		fmt.Println("Error creating tun device:", err)
		return
	}
	defer ifce.Close()

	err = netTool.SetTunUp(ifce.Name())
	if err != nil {
		fmt.Println("Error setting tun up:", err)
		return
	}

	fmt.Println("Listening on udp 8299 success....")

	remoteMap := make(map[string]*net.UDPAddr)

	fmt.Println("Listening on udp 8299")
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 用 UdpPackage 读取数据
	p := protocol.UdpPackage{}
	for {
		// 读取数据
		fmt.Println("Reading package....")
		remoteUdpAddr, err := p.ReadPackage(conn)
		if err != nil {
			panic(err)
		}

		// p.Payload 是一个ip包，  解析出ip包的源ip和目的ip， 打印出来
		ipHeader, err := ipv4.ParseHeader(p.Payload)
		if err != nil {
			panic(err)
		}

		fmt.Println(ipHeader.Checksum)
		fmt.Println(checksum(p.Payload[0:20]))

		// 获取本机ip
		ifceLocal, err := net.InterfaceByName("eth0")
		if err != nil {
			panic(err)
		}

		addrs, err := ifceLocal.Addrs()
		if err != nil {
			panic(err)
		}

		var localIp string
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				panic(err)
			}
			if ip.To4() != nil {
				localIp = ip.String()
				break
			}
		}

		fmt.Println(localIp)

		// 修改 p.Payload 源ip
		copy(p.Payload[12:16], net.ParseIP(localIp).To4())
		// 重新计算ip包 crc校验和
		sum := checksum(p.Payload[0:20])
		p.Payload[10] = byte(sum >> 8)
		p.Payload[11] = byte(sum)

		if _, ok := remoteMap[p.SrcIp]; !ok {
			remoteMap[p.SrcIp] = remoteUdpAddr
		}

		fmt.Println(binary.BigEndian.Uint16(p.Payload[10:12]))
		fmt.Println(sum)

		size, err := ifce.Write(p.Payload)
		if err != nil {
			fmt.Println("Error writing to tun device:", err)
			return
		} else {
			fmt.Println("Write to tun device:", size)
		}

		resp := make([]byte, 1500)
		readSize, err := ifce.Read(resp)
		if err != nil {
			fmt.Println("Error reading from tun device:", err)
			return
		} else {
			fmt.Println("Read from tun device:", readSize)
		}
	}
}

func checksum(header []byte) uint16 {
	// 检验和置0
	header[10] = 0
	header[11] = 0

	var sum uint32

	for i := 0; i < len(header); i += 2 {
		sum += uint32(header[i])<<8 | uint32(header[i+1])
	}
	// 将进位加到低 16 位上
	for sum>>16 > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	// 取反得到校验和
	return uint16(^sum)
}
