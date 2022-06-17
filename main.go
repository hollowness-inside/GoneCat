package main

import (
	"fmt"
	"net"
	"os"
)

func help() {
	println("Use: gnc [options] address:port")
	println("\t-u - Use UDP connection")
	println("\t-t - Use TCP connection (Default)")
	println("\t-C - Send CRLF as line-ending (Default is none)")
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
		case "-C":
			args.SendCRLF = true
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
