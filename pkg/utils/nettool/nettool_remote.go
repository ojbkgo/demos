package nettool

import (
	"fmt"
	"os/exec"

	"github.com/songgao/water"
)

type linuxNetTool struct{}

func NewLinuxNetTool() *linuxNetTool {
	return &linuxNetTool{}
}

func (l *linuxNetTool) CreateTun(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	ifce, err := water.New(config)
	if err != nil {
		fmt.Println("Error creating tun device:", err)
		return nil, err
	}

	return ifce, nil
}

func (l *linuxNetTool) SetTunNetIP(name string, ip, mask string) error {
	out, err := exec.Command("ifconfig", name, "inet", ip, "netmask", mask).CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting tun net IP: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}

func (l *linuxNetTool) SetTunUp(name string) error {
	out, err := exec.Command("ifconfig", name, "up").CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting tun up: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}

func (l *linuxNetTool) AddRoute(net string, gw string) error {
	_, err := exec.Command("route", "add", "-net", net, "gw", gw).Output()
	if err != nil {
		fmt.Println("Error adding route:", err)
		return err
	}
	return nil
}

func (l *linuxNetTool) DelRoute(net string, gw string) error {
	_, err := exec.Command("route", "del", "-net", net, "gw", gw).Output()
	if err != nil {
		fmt.Println("Error deleting route:", err)
		return err
	}
	return nil
}

func (l *linuxNetTool) UninstallTun(name string) error {
	// Bring down the interface
	out, err := exec.Command("ifconfig", name, "down").CombinedOutput()
	if err != nil {
		fmt.Printf("Error bringing down interface: %v, output: %s\n", err, string(out))
		return err
	}

	// Delete the interface
	out, err = exec.Command("ip", "tuntap", "del", "mode", "tun", "dev", name).CombinedOutput()
	if err != nil {
		fmt.Printf("Error deleting interface: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}

func (l *linuxNetTool) AddSnat(srcNet string) error {
	_, err := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", srcNet, "-j", "MASQUERADE").Output()
	if err != nil {
		fmt.Println("Error adding snat:", err)
		return err
	}
	return nil
}

func (l *linuxNetTool) DelSnat(srcNet string) error {
	_, err := exec.Command("iptables", "-t", "nat", "-D", "POSTROUTING", "-s", srcNet, "-j", "MASQUERADE").Output()
	if err != nil {
		fmt.Println("Error deleting snat:", err)
		return err
	}
	return nil
}
