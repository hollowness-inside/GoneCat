package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Use: gnc [options] address port")
	fmt.Println("\t-4 - Use IPv4")
	fmt.Println("\t-6 - Use IPv6")
	fmt.Println("\t-C - Send CRLF as line-ending")
	fmt.Println("\t-d - Detach from stdin")
	fmt.Println("\t-u - UDP mode")
}

func main() {
	gct := GoneCat{}
	gct.UseDefaults()

	arg := 1

	for arg < len(os.Args) {
		cur := os.Args[arg]
		switch cur {
		case "-4":
			gct.OnlyIpv4 = true
		case "-6":
			gct.OnlyIpv6 = true
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
			if arg+1 >= len(os.Args) {
				help()
				return
			}

			gct.AddrStr = cur
			gct.AddrPort = os.Args[arg+1]
			arg += 1
		}

		arg++
	}

	err := gct.Execute()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
