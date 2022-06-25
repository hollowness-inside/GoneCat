package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type GoneCat struct {
	AddrStr    string
	AddrPort   string
	Address    *net.TCPAddr
	Network    string
	Listening  bool
	OnlyIpv4   bool
	OnlyIpv6   bool
	Tcp        bool
	SendCRLF   bool
	ReadStdin  bool
	ReadPipe   bool
	BufferSize int
	Addr       net.TCPAddr
}

func (gc *GoneCat) UseDefaults() {
	gc.Listening = false
	gc.OnlyIpv4 = false
	gc.OnlyIpv6 = false
	gc.Tcp = true
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

func (gc *GoneCat) resolveAddress() {
	only := ""

	if gc.OnlyIpv4 {
		only = "4"
	}

	if gc.OnlyIpv6 {
		only = "6"
	}

	var protocol string
	if gc.Tcp {
		protocol = "tcp"
	} else {
		protocol = "udp"
	}

	gc.Network = protocol + only

	ip := net.ParseIP(gc.AddrStr)
	port, err := strconv.Atoi(gc.AddrPort)
	if err != nil {
		panic("The given port is not a number")
	}

	gc.Address = &net.TCPAddr{IP: ip, Port: port, Zone: ""}
}

func (gc *GoneCat) doListen() error {
	listener, err := net.Listen(gc.Network, gc.Address.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go gc.handleConnection(conn)
	}
}

func (gc *GoneCat) doConnect() error {
	conn, err := net.Dial(gc.Network, gc.Address.String())
	if err != nil {
		return err
	}

	gc.handleConnection(conn)
	return nil
}

func (gc *GoneCat) handleConnection(conn net.Conn) {
	defer conn.Close()

	if gc.ReadPipe {
		go gc.streamPipe(conn)
	}

	if gc.ReadStdin {
		go gc.streamStdin(conn)
	}

	io.Copy(os.Stdout, conn)
}

func (gc *GoneCat) streamPipe(conn net.Conn) {
	for {
		_, err := io.CopyN(conn, os.Stdin, int64(gc.BufferSize))
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
}

func (gc *GoneCat) streamStdin(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()

		if gc.SendCRLF {
			fmt.Fprintln(conn, str)
		} else {
			conn.Write([]byte(str))
		}
	}
}
