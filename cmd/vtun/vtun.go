package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/net/ipv4"

	"github.com/ojbkgo/demos/pkg/protocol"
	"github.com/ojbkgo/demos/pkg/utils/nettool"
)

var remoteAddr = flag.String("remote", "", "remote address")

const tunName = "utun99"

func main() {

	flag.Parse()
	if *remoteAddr == "" {
		fmt.Println("remote address is required")
		return
	}

	macNetTool := nettool.NewMacNetTool()
	var err error
	ifce, err := macNetTool.CreateTun(tunName)
	if err != nil {
		fmt.Println("Error creating tun device:", err)
		return
	}

	defer ifce.Close()
	err = macNetTool.SetTunNetIP(ifce.Name(), "192.168.0.99", "255.255.255.0")

	if err != nil {
		fmt.Println("Error setting tun net IP:", err)
		return
	}
	err = macNetTool.SetTunUp(tunName)
	if err != nil {
		fmt.Println("Error setting tun up:", err)
		return
	}
	err = macNetTool.AddRoute("172.24.224.0/20", "192.168.0.100")
	if err != nil {
		fmt.Println("Error adding route:", err)
		return
	}

	fmt.Println(tunName + " created, dial:" + *remoteAddr + " success....")
	remoteUdpAddr, err := net.ResolveUDPAddr("udp", *remoteAddr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}
	// udp 打开 远程 8299 端口
	conn, err := net.DialUDP("udp", nil, remoteUdpAddr)
	if err != nil {
		fmt.Println("Error dialing udp:", err)
		return
	}

	fmt.Println("dial udp success....")

	go func() {

		for {
			buffer := make([]byte, 1500)
			readSize, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading from udp:", err)
			}

			fmt.Println("Read from udp success.... size:", readSize)
			_, err = ifce.Write(buffer[:readSize])
			if err != nil {
				fmt.Println("Error writing to tun device:", err)
			} else {
				fmt.Println("Write to tun device success....")
			}
		}

	}()

	for {
		buffer := make([]byte, 1500)
		n, err := ifce.Read(buffer)

		data := protocol.NewUdpPackage("192.168.0.99", "172.0.0.1", buffer[:n])
		if err != nil {
			fmt.Println("Error reading from tun device:", err)
			return
		}

		ipHeader, err := ipv4.ParseHeader(data.Payload)
		if err != nil {
			fmt.Println("Error parsing ip header:", err)
			return
		}
		fmt.Println("ip package src ip:", ipHeader.Src.String())
		fmt.Println("ip package dest ip:", ipHeader.Dst.String())

		fmt.Println("Read from tun device:", len(buffer[:n]))

		size, err := conn.Write(data.Serialize())
		if err != nil {
			fmt.Println("Error writing to udp:", err)
			return
		}

		fmt.Println("Write to udp success.... size:", size)

	}
}
