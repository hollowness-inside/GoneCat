package gonecat

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func (gc *GoneCat) tcpListen() error {
	listener, err := net.Listen(gc.Network, gc.Address.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go gc.handleTCP(conn)
	}
}

func (gc *GoneCat) tcpConnect() error {
	conn, err := net.Dial(gc.Network, gc.Address.String())
	if err != nil {
		return err
	}

	gc.handleTCP(conn)
	return nil
}

func (gc *GoneCat) handleTCP(conn net.Conn) {
	defer conn.Close()

	if gc.ReadPipe {
		go gc.streamPipeTCP(conn)
	}

	if gc.ReadStdin {
		go gc.streamStdinTCP(conn)
	}

	io.Copy(os.Stdout, conn)
}

func (gc *GoneCat) streamPipeTCP(conn net.Conn) {
	for {
		_, err := io.CopyN(conn, os.Stdin, int64(gc.BufferSize))
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
}

func (gc *GoneCat) streamStdinTCP(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()

		if gc.SendCRLF {
			fmt.Fprintln(conn, str)
		} else {
			conn.Write([]byte(str))
		}
	}
}
