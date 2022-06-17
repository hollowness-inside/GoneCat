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
	println("\t-d - Do not attempt to read from stdin")
	println("\t-4 - Use only IPv4")
	println("\t-6 - Use only IPv6")
}

func main() {
	gct := GoneCat{}
	gct.UseDefaults()

	arg := 1

	for arg < len(os.Args) {
		cur := os.Args[arg]
		switch cur {
		case "-4":
			gct.Ipv4 = true
		case "-6":
			gct.Ipv6 = true
		case "-u":
			gct.Tcp = false
		case "-C":
			gct.SendCRLF = true
		case "-l":
			gct.Listening = true
		case "-d":
			gct.ReadStdin = false
		case "-h", "--help":
			help()
			return
		default:
			addr, err := net.ResolveTCPAddr("", cur)
			if err != nil {
				help()
				return
			}
			gct.Addr = *addr
		}

		arg++
	}

	err := gct.Execute()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
