package main

import (
	"bufio"
	"errors"
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

	if gc.Protocol == "tcp" {
		if gc.Listening {
			return gc.tcpListen()
		}
		return gc.tcpConnect()
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

func (gc *GoneCat) tcpListen() error {
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

func (gc *GoneCat) tcpConnect() error {
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
