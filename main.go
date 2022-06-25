package main

import (
	"fmt"
	"joshua/green/gonecat/gonecat"
	"os"
)

func help() {
	fmt.Println("Use: gnc [options] address port")
	fmt.Println("\t-4 - Use IPv4")
	fmt.Println("\t-6 - Use IPv6")
	fmt.Println("\t-C - Send CRLF as line-ending")
	fmt.Println("\t-d - Detach from stdin")
	fmt.Println("\t-u - Use UDP")
}

func main() {
	gct := gonecat.GoneCatArguments{}
	gct.UseDefaults()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		gct.ReadPipe = true
	}

	arg := 1

	for arg < len(os.Args) {
		cur := os.Args[arg]
		switch cur {
		case "-4":
			gct.IPVersion = 4
		case "-6":
			gct.IPVersion = 6
		case "-C":
			gct.SendCRLF = true
		case "-u":
			gct.Protocol = "udp"
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

	gonecat := gonecat.GetCat(gct)
	if gonecat == nil {
		panic("an error occured on trying to run gonecat")
	}

	err = gonecat.Execute()
	if err != nil {
		panic(err)
	}
}
