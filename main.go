package main

import (
	"fmt"
	"joshua/green/gonecat/gonecat"
	"log"
	"os"
)

const HelpMsg = `Usage: gnc [options] address port
	-4	Use IPv4
	-6	Use IPv6
	-C	Do not send CRLF as line-ending
	-d	Detach from stdin
	-l	Listen
	-o file	Redirect output to file
	-u	Use UDP
`

func main() {
	gct := ParseArguments()
	if gct == nil {
		fmt.Println(HelpMsg)
		return
	}

	gonecat := gonecat.GetCat(gct)
	err := gonecat.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func ParseArguments() *gonecat.GCArguments {
	if len(os.Args) == 1 {
		return nil
	}

	gct := new(gonecat.GCArguments)
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
		case "-o":
			i++
			path := os.Args[i]
			output, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}

			gct.Output = output
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
		log.Fatal(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}
