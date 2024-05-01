package gonecat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// TcpCat is a TCP client/server that reads from stdin and writes to stdout.
type TcpCat struct {
	*GCArguments
	Address *net.TCPAddr
}

// Execute runs the TCP client/server.
func (tc *TcpCat) Execute() error {
	if tc.Listening {
		return tc.listen()
	}

	return tc.connect()
}

// listen listens for incoming connections and handles them.
func (tc *TcpCat) listen() error {
	listener, err := net.ListenTCP(tc.Network, tc.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go tc.handle(&GCCon{conn})
	}
}

// connect connects to the remote server and handles the connection.
func (tc *TcpCat) connect() error {
	conn, err := net.DialTCP(tc.Network, nil, tc.Address)
	if err != nil {
		return err
	}

	tc.handle(&GCCon{conn})
	return nil
}

// handle handles the connection.
func (tc *TcpCat) handle(conn *GCCon) {
	defer conn.Close()
	defer tc.Output.Close()

	if tc.ReadPipe {
		tc.streamPipe(conn)
	}

	if tc.ReadStdin {
		go tc.streamStdin(conn)
	}

	io.Copy(tc.Output, conn)
}

// streamPipe streams the pipe to the connection.
func (tc *TcpCat) streamPipe(conn *GCCon) {
	for {
		_, err := io.CopyN(conn, os.Stdin, int64(tc.BufferSize))
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// streamStdin streams stdin to the connection.
func (tc *TcpCat) streamStdin(conn *GCCon) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()

		if tc.SendCRLF {
			fmt.Fprintln(conn, str)
		} else {
			conn.Write([]byte(str))
		}
	}
}
