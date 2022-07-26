package main

import (
	"fmt"
	"joshua/green/gonecat/gonecat"
	"log"
	"os"
	"strconv"
)

const HelpMsg = `Usage: gnc [options] address port
	-4		Use IPv4
	-6		Use IPv6
	-C		Do not send CRLF as line-ending
	-d		Detach from stdin
	-h		This help text
	-I length	TCP receive buffer length
	-l		Listen mode
	-o file		Redirect output to file	
	-u		UDP mode
`

func main() {
	gct := ParseArguments()
	if gct == nil {
		fmt.Print(HelpMsg)
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
		case "-d":
			gct.ReadStdin = false
		case "-h", "--help":
			return nil
		case "-I":
			i++
			bsize, err := strconv.Atoi(os.Args[i])
			if err != nil {
				return nil
			}

			gct.BufferSize = bsize
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
		default:
			if i+1 >= len(os.Args) {
				return nil
			}

			gct.AddrStr = arg
			gct.AddrPort = os.Args[i+1]
			i++
		}

		i++
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
