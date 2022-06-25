package gonecat

import (
	"errors"
	"net"
	"strconv"
)

type GoneCat struct {
	AddrStr    string
	AddrPort   string
	Address    *net.TCPAddr
	Network    string
	Listening  bool
	IPVersion  uint8
	Protocol   string
	SendCRLF   bool
	ReadStdin  bool
	ReadPipe   bool
	BufferSize int
	Addr       net.TCPAddr
}

func (gc *GoneCat) UseDefaults() {
	gc.Listening = false
	gc.IPVersion = 0
	gc.Protocol = "tcp"
	gc.SendCRLF = false
	gc.ReadStdin = true
	gc.ReadPipe = false
	gc.BufferSize = 1024
}

func (gc *GoneCat) Execute() error {
	gc.resolveAddress()

	if gc.Listening {
		return gc.doListen()
	}
	return gc.doConnect()
}

func (gc *GoneCat) doListen() error {
	if gc.Protocol == "tcp" {
		return gc.tcpListen()
	} else if gc.Protocol == "udp" {
		return gc.udpListen()
	}

	return errors.New("cannot listen: no protocol provided")
}

func (gc *GoneCat) doConnect() error {
	if gc.Protocol == "tcp" {
		return gc.tcpConnect()
	} else if gc.Protocol == "udp" {
		return gc.udpConnect()
	}

	return errors.New("cannot connect: no protocol provided")
}

func (gc *GoneCat) resolveAddress() {
	var version string = ""
	switch gc.IPVersion {
	case 4:
		version = "4"
	case 6:
		version = "6"
	}

	gc.Network = gc.Protocol + version

	ip := net.ParseIP(gc.AddrStr)
	port, err := strconv.Atoi(gc.AddrPort)
	if err != nil {
		panic("The given port is not a number")
	}

	gc.Address = &net.TCPAddr{IP: ip, Port: port, Zone: ""}
}
