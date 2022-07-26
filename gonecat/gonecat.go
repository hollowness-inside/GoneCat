package gonecat

import (
	"log"
	"net"
	"strconv"
)

type GCCon struct {
	net.Conn
}

type GCArguments struct {
	AddrStr    string
	AddrPort   string
	Network    string
	Listening  bool
	Protocol   string
	IPVersion  string
	SendCRLF   bool
	ReadStdin  bool
	ReadPipe   bool
	BufferSize int
}

type GoneCat interface {
	Execute() error
	listen() error
	connect() error
	handle(conn *GCCon)
	streamPipe(conn *GCCon)
	streamStdin(conn *GCCon)
}

func (gc *GCArguments) UseDefaults() {
	gc.Listening = false
	gc.IPVersion = ""
	gc.Protocol = "tcp"
	gc.SendCRLF = true
	gc.ReadStdin = true
	gc.ReadPipe = false
	gc.BufferSize = 1024
}

func GetCat(gc *GCArguments) GoneCat {
	addr := gc.resolveAddress()

	switch gc.Protocol {
	case "tcp":
		return &TcpCat{gc, addr.(*net.TCPAddr)}
	case "udp":
		return &UdpCat{gc, addr.(*net.UDPAddr)}
	default:
		log.Fatalf("Wrong protocol name %s", gc.Protocol)
		return nil
	}
}

func (gc *GCArguments) resolveAddress() net.Addr {
	gc.Network = gc.Protocol + gc.IPVersion

	ip := net.ParseIP(gc.AddrStr)
	port, err := strconv.Atoi(gc.AddrPort)
	if err != nil {
		log.Fatal("The given port is not a number")
	}

	if gc.Protocol == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port, Zone: ""}
	} else if gc.Protocol == "udp" {
		return &net.UDPAddr{IP: ip, Port: port, Zone: ""}
	} else {
		return nil
	}
}
