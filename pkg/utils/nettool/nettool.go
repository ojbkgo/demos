package nettool

import (
	"github.com/songgao/water"
)

type NetTool interface {
	// CreateTun create tun device
	CreateTun(name string) (*water.Interface, error)
	// SetTunNetIP set tun net ip, format like 192.168.x.x/24
	SetTunNetIP(name string, ip, mask string) error
	// SetTunUp set tun up
	SetTunUp(name string) error
	// AddRoute add route
	AddRoute(net string, gw string) error
	// DelRoute del route
	DelRoute(net string, gw string) error
	// UninstallTun uninstall tun
	UninstallTun(name string) error
}
