package main

import (
	"net"
	"os"
)

type Arguments struct {
	Listening bool
	Ipv4      bool
	Ipv6      bool
	Udp       bool
	Address   net.TCPAddr
}

func main() {
	args := Arguments{}
	arg := 1

	for arg < len(os.Args) {
		cur := os.Args[arg]
		switch cur {
		case "-4":
			args.Ipv4 = true
		case "-6":
			args.Ipv6 = true
		case "-u":
			args.Udp = true
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
}

func help() {
	println("Use: nc [-46ul] address:port")
}
