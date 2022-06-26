package main

import (
	"fmt"
	"joshua/green/gonecat/gonecat"
	"os"
)

func ShowHelp() {
	fmt.Println("Use: gnc [options] address port")
	fmt.Println("\t-4 - Use IPv4")
	fmt.Println("\t-6 - Use IPv6")
	fmt.Println("\t-C - Do not send CRLF as line-ending")
	fmt.Println("\t-d - Detach from stdin")
	fmt.Println("\t-u - Use UDP")
}

func main() {
	gct := ParseArguments()
	if gct == nil {
		ShowHelp()
		return
	}

	gonecat := gonecat.GetCat(gct)
	if gonecat == nil {
		panic("an error occured on trying to get gonecat")
	}

	err := gonecat.Execute()
	if err != nil {
		panic(err)
	}
}

func ParseArguments() *gonecat.GoneCatArguments {
	if len(os.Args) == 1 {
		return nil
	}

	gct := new(gonecat.GoneCatArguments)
	gct.UseDefaults()

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-4":
			gct.IPVersion = "4"
		case "-6":
			gct.IPVersion = "6"
		case "-C":
			gct.SendCRLF = false
		case "-u":
			gct.Protocol = "udp"
		case "-l":
			gct.Listening = true
		case "-d":
			gct.ReadStdin = false
		case "-h", "--help":
			return nil
		default:
			if i+1 >= len(os.Args) {
				return nil
			}

			gct.AddrStr = arg
			gct.AddrPort = os.Args[i+1]
			i += 1
		}

		i += 1
	}

	if gct.AddrStr == "" || gct.AddrPort == "" {
		return nil
	}

	gct.ReadPipe = isPipeConnected()

	return gct
}

func isPipeConnected() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}
