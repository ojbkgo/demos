package nettool

import (
	"fmt"
	"os/exec"

	"github.com/songgao/water"
)

type macNetTool struct {
}

func NewMacNetTool() *macNetTool {
	return &macNetTool{}
}

func (m *macNetTool) CreateTun(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	ifce, err := water.New(config)
	if err != nil {
		fmt.Println("Error creating tun device:", err)
		return nil, err
	}

	// route add -net
	return ifce, nil
}

func (m *macNetTool) SetTunNetIP(name string, ip, mask string) error {
	out, err := exec.Command("ifconfig", name, "inet", ip, "192.168.0.100", "netmask", mask).CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting tun net IP: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}

func (m *macNetTool) SetTunUp(name string) error {
	out, err := exec.Command("ifconfig", name, "up").CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting tun up: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}

func (m *macNetTool) AddRoute(net string, gw string) error {
	_, err := exec.Command("route", "add", "-net", net, gw).Output()
	if err != nil {
		fmt.Println("Error adding route:", err)
		return err
	}
	return nil
}

func (m *macNetTool) DelRoute(net string, gw string) error {
	_, err := exec.Command("route", "delete", "-net", net, gw).Output()
	if err != nil {
		fmt.Println("Error deleting route:", err)
		return err
	}
	return nil
}

func (m *macNetTool) UninstallTun(name string) error {
	// Bring down the interface
	out, err := exec.Command("ifconfig", name, "down").CombinedOutput()
	if err != nil {
		fmt.Printf("Error bringing down interface: %v, output: %s\n", err, string(out))
		return err
	}

	// Delete the interface
	out, err = exec.Command("ifconfig", name, "destroy").CombinedOutput()
	if err != nil {
		fmt.Printf("Error deleting interface: %v, output: %s\n", err, string(out))
		return err
	}

	return nil
}
