package gonecat

import (
	"net"
	"strconv"
)

type GCCon struct {
	net.Conn
}

type GoneCat interface {
	Execute() error
	listen() error
	connect() error
	handle(conn *GCCon)
	streamPipe(conn *GCCon)
	streamStdin(conn *GCCon)
}

type GoneCatArguments struct {
	AddrStr    string
	AddrPort   string
	Network    string
	Listening  bool
	IPVersion  uint8
	Protocol   string
	SendCRLF   bool
	ReadStdin  bool
	ReadPipe   bool
	BufferSize int
}

func (gc *GoneCatArguments) UseDefaults() {
	gc.Listening = false
	gc.IPVersion = 0
	gc.Protocol = "tcp"
	gc.SendCRLF = true
	gc.ReadStdin = true
	gc.ReadPipe = false
	gc.BufferSize = 1024
}

func GetCat(gc GoneCatArguments) GoneCat {
	gc.resolveAddress()

	addr := gc.resolveAddress()

	if gc.Protocol == "tcp" {
		return TcpCat{gc, addr.(*net.TCPAddr)}
	} else if gc.Protocol == "udp" {
		return UdpCat{gc, addr.(*net.UDPAddr)}
	}

	return nil
}

func (gc *GoneCatArguments) resolveAddress() net.Addr {
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

	if gc.Protocol == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port, Zone: ""}
	} else if gc.Protocol == "udp" {
		return &net.UDPAddr{IP: ip, Port: port, Zone: ""}
	} else {
		return nil
	}
}
