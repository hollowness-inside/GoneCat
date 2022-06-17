package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

type GoneCat struct {
	Listening bool
	Ipv4      bool
	Ipv6      bool
	Tcp       bool
	Addr      net.TCPAddr
}

func (gc *GoneCat) UseDefaults() {
	gc.Listening = false
	gc.Ipv4 = true
	gc.Ipv6 = false
	gc.Tcp = true
	gc.Addr = net.TCPAddr{}
}

func (gc *GoneCat) Network() string {
	if gc.Tcp {
		return "tcp"
	} else {
		return "udp"
	}
}

func (gc *GoneCat) Address() string {
	return gc.Addr.String()
}

func (gc *GoneCat) Execute() error {
	if gc.Listening {
		return gc.doListen()
	} else {
		return gc.doConnect()
	}
}

func (gc *GoneCat) doListen() error {
	listener, err := net.Listen(gc.Network(), gc.Address())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func (gc *GoneCat) doConnect() error {
	conn, err := net.Dial(gc.Network(), gc.Address())
	if err != nil {
		return err
	}

	handleConnection(conn)
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() {
			str := scanner.Text()
			str = strings.TrimRight(str, "\r\n")
			conn.Write([]byte(str))
		}
	}()

	io.Copy(os.Stdout, conn)
}
