package gonecat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func (gc *GoneCat) udpListen() error {
	conn, err := net.ListenUDP(gc.Network, (*net.UDPAddr)(gc.Address))
	if err != nil {
		return err
	}
	defer conn.Close()

	gc.handleUDP(conn)
	return nil
}

func (gc *GoneCat) udpConnect() error {
	conn, err := net.DialUDP(gc.Network, nil, (*net.UDPAddr)(gc.Address))
	if err != nil {
		return err
	}
	defer conn.Close()

	gc.handleUDP(conn)
	return nil
}

func (gc *GoneCat) handleUDP(conn *net.UDPConn) {
	if gc.ReadPipe {
		go gc.streamPipeUDP(conn)
	}

	if !gc.Listening && gc.ReadStdin {
		go gc.streamStdinUDP(conn)
	}

	io.Copy(os.Stdout, conn)
}

func (gc *GoneCat) streamPipeUDP(conn *net.UDPConn) {
	for {
		buffer := make([]byte, gc.BufferSize)
		_, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		conn.Write(buffer)
	}
}

func (gc *GoneCat) streamStdinUDP(conn *net.UDPConn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()

		if gc.SendCRLF {
			conn.Write([]byte(fmt.Sprintln(str)))
		} else {
			conn.Write([]byte(str))
			// conn.WriteMsgUDP([]byte(str), nil, nil)
		}
	}
}
