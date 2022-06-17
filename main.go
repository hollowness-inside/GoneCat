package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Arguments struct {
	Listening bool
	Ipv4      bool
	Ipv6      bool
	Tcp       bool
	Addr      net.TCPAddr
}

func (a Arguments) Default() Arguments {
	a.Listening = false
	a.Ipv4 = true
	a.Ipv6 = false
	a.Tcp = true
	a.Addr = net.TCPAddr{}

	return a
}

func (a *Arguments) Network() string {
	if a.Tcp {
		return "tcp"
	} else {
		return "udp"
	}
}

func (a *Arguments) Address() string {
	return a.Addr.String()
}

func main() {
	args := Arguments{}.Default()
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

	err := Execute(args)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func help() {
	println("Use: nc [-46ul] address:port")
}

func Execute(args Arguments) error {
	if args.Listening {
		return doListen(args)
	} else {
		return doConnect(args)
	}
}

func doListen(args Arguments) error {
	listener, err := net.Listen(args.Network(), args.Address())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func doConnect(args Arguments) error {
	conn, err := net.Dial(args.Network(), args.Address())
	if err != nil {
		return err
	}

	handleConnection(conn)
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			var str string
			fmt.Scanln(&str)
			str = strings.TrimRight(str, "\r\n")
			conn.Write([]byte(str))
		}
	}()

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatalf("Cannot copy from Conn into Stdout. Conn: %s", conn.RemoteAddr().String())
		}
		// fmt.Printf("%s has been disconnected\n", conn.RemoteAddr().String())
	}()

	for {
	}
}
