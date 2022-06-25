package gonecat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type UdpCat struct {
	GoneCatArguments
	Address *net.UDPAddr
}

func (uc UdpCat) Execute() error {
	if uc.Listening {
		return uc.listen()
	}

	return uc.connect()
}

func (uc UdpCat) listen() error {
	conn, err := net.ListenUDP(uc.Network, uc.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	uc.handle(&GCCon{conn})
	return nil
}

func (uc UdpCat) connect() error {
	conn, err := net.DialUDP(uc.Network, nil, (*net.UDPAddr)(uc.Address))
	if err != nil {
		return err
	}
	defer conn.Close()

	uc.handle(&GCCon{conn})
	return nil
}

func (uc UdpCat) handle(conn *GCCon) {
	if uc.ReadPipe {
		go uc.streamPipe(conn)
	}

	if !uc.Listening && uc.ReadStdin {
		go uc.streamStdin(conn)
	}

	io.Copy(os.Stdout, conn)
}

func (uc UdpCat) streamPipe(conn *GCCon) {
	for {
		buffer := make([]byte, uc.BufferSize)
		_, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		conn.Write(buffer)
	}
}

func (uc UdpCat) streamStdin(conn *GCCon) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()

		if uc.SendCRLF {
			conn.Write([]byte(fmt.Sprintln(str)))
		} else {
			conn.Write([]byte(str))
		}
	}
}
