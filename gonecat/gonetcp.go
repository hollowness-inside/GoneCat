package gonecat

import "net"

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

		go gc.handleConnection(conn)
	}
}

func (gc *GoneCat) tcpConnect() error {
	conn, err := net.Dial(gc.Network, gc.Address.String())
	if err != nil {
		return err
	}

	gc.handleConnection(conn)
	return nil
}
