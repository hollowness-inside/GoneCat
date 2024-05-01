package gonecat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// UdpCat is a struct that holds the arguments for the UDP cat command
type UdpCat struct {
	*GCArguments
	Address *net.UDPAddr
}

// Execute runs the UDP client/server.
func (uc *UdpCat) Execute() error {
	if uc.Listening {
		return uc.listen()
	}

	return uc.connect()
}

// listen listens for incoming connections and handles them.
func (uc *UdpCat) listen() error {
	conn, err := net.ListenUDP(uc.Network, uc.Address)
	if err != nil {
		return err
	}

	uc.handle(&GCCon{conn})
	return nil
}

// connect connects to the server and handles the connection.
func (uc *UdpCat) connect() error {
	conn, err := net.DialUDP(uc.Network, nil, (*net.UDPAddr)(uc.Address))
	if err != nil {
		return err
	}

	uc.handle(&GCCon{conn})
	return nil
}

// handle handles the connection.
func (uc *UdpCat) handle(conn *GCCon) {
	defer conn.Close()
	defer uc.Output.Close()

	if uc.ReadPipe {
		uc.streamPipe(conn)
	}

	if !uc.Listening && uc.ReadStdin {
		go uc.streamStdin(conn)
	}

	io.Copy(uc.Output, conn)
}

// streamPipe streams the pipe to the connection.
func (uc *UdpCat) streamPipe(conn *GCCon) {
	for {
		buffer := make([]byte, uc.BufferSize)
		_, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		conn.Write(buffer)
	}
}

// streamStdin streams stdin to the connection.
func (uc *UdpCat) streamStdin(conn *GCCon) {
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
