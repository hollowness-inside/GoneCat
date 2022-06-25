package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

type GoneCat struct {
	AddrStr   string
	AddrPort  string
	Address   string
	Network   string
	Listening bool
	OnlyIpv4  bool
	OnlyIpv6  bool
	Tcp       bool
	SendCRLF  bool
	ReadStdin bool
	ReadPipe  bool
	Addr      net.TCPAddr
}

func (gc *GoneCat) UseDefaults() {
	gc.Listening = false
	gc.OnlyIpv4 = false
	gc.OnlyIpv6 = false
	gc.Tcp = true
	gc.SendCRLF = false
	gc.ReadStdin = true
	gc.ReadPipe = false
}

func (gc *GoneCat) Execute() error {
	gc.resolveAddress()
	gc.resolveNetwork()

	if gc.Listening {
		return gc.doListen()
	}

	return gc.doConnect()
}

func (gc *GoneCat) resolveAddress() {
	ipv4, _ := regexp.MatchString(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`, gc.AddrStr)
	if ipv4 {
		gc.Address = gc.AddrStr + ":" + gc.AddrPort
	} else {
		gc.Address = "[" + gc.AddrStr + "]:" + gc.AddrPort
	}
}

func (gc *GoneCat) resolveNetwork() {
	only := ""

	if gc.OnlyIpv4 {
		only = "4"
	}

	if gc.OnlyIpv6 {
		only = "6"
	}

	var protocol string
	if gc.Tcp {
		protocol = "tcp"
	} else {
		protocol = "udp"
	}

	gc.Network = protocol + only
}

func (gc *GoneCat) doListen() error {
	listener, err := net.Listen(gc.Network, gc.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go gc.handleConnection(conn)
	}
}

func (gc *GoneCat) doConnect() error {
	conn, err := net.Dial(gc.Network, gc.Address)
	if err != nil {
		return err
	}

	gc.handleConnection(conn)
	return nil
}

func (gc *GoneCat) handleConnection(conn net.Conn) {
	defer conn.Close()

	if gc.ReadPipe {
		reader := bufio.NewReader(os.Stdin)

		go func() {
			var data []byte
			for {
				b, err := reader.ReadByte()
				if err == io.EOF {
					break
				}
				data = append(data, b)
			}

			conn.Write(data)
			data = nil
		}()
	} else {
		if gc.ReadStdin {
			scanner := bufio.NewScanner(os.Stdin)

			go func() {
				for scanner.Scan() {
					str := scanner.Text()

					if gc.SendCRLF {
						fmt.Fprintln(conn, str)
					} else {
						conn.Write([]byte(str))
					}
				}
			}()
		}
	}

	io.Copy(os.Stdout, conn)

}
