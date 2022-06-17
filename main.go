package main

import (
	"bufio"
	"fmt"
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

func help() {
	println("Use: gnc [options] address:port")
	println("\t-u - Use UDP connection")
	println("\t-t - Use TCP connection (Default)")
	println("\t-4 - Use only IPv4")
	println("\t-6 - Use only IPv6")
}

func main() {
	args := GoneCat{}
	args.UseDefaults()

	arg := 1

	for arg < len(os.Args) {
		cur := os.Args[arg]
		switch cur {
		case "-4":
			args.Ipv4 = true
		case "-6":
			args.Ipv6 = true
		case "-u":
			args.Tcp = false
		case "-l":
			args.Listening = true
		case "-h", "--help":
			help()
			return
		default:
			addr, err := net.ResolveTCPAddr("", cur)
			if err != nil {
				help()
				return
			}
			args.Addr = *addr
		}

		arg++
	}

	err := args.Execute()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
