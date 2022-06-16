package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Arguments struct {
	Listening bool
	Ipv4      bool
	Ipv6      bool
	Tcp       bool
	Address   net.TCPAddr
}

func main() {
	args := Arguments{}
	args.Listening = false
	args.Ipv4 = true
	args.Ipv6 = false
	args.Tcp = true
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
			args.Address = *addr
		}

		arg++
	}

	err := execute(args)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func help() {
	println("Use: nc [-46ul] address:port")
}

func execute(args Arguments) error {
	if args.Listening {
		return DoListen(args)
	} else {
		return DoConnect(args)
	}
}

func DoListen(args Arguments) error {
	return nil
}

func DoConnect(args Arguments) error {
	var network string
	if args.Tcp {
		network = "tcp"
	} else {
		network = "udp"
	}

	conn, err := net.Dial(network, args.Address.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		io.Copy(os.Stdout, os.Stdin)
		panic("Cannot copy from Stdin into Stdout")
	}()

	go func() {
		io.Copy(conn, os.Stdin)
		panic("Cannot copy from Stdin into Conn")
	}()

	go func() {
		io.Copy(os.Stdout, conn)
		panic("Cannot copy from Conn into Stdout")
	}()

	for {
	}
}
