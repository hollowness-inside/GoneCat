package gonecat

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

// GCCon is a wrapper for net.Conn
type GCCon struct {
	net.Conn
}

// GoneCat is an interface representing a connection to a remote process.
type GoneCat interface {
	// Execute runs the command on the remote process.
	Execute() error

	// Listen starts listening for incoming connections.
	listen() error

	// Connect establishes a connection to the remote process.
	connect() error

	// Handle processes an incoming connection.
	handle(conn *GCCon)

	// StreamPipe streams data from the remote process to the local process.
	streamPipe(conn *GCCon)

	// StreamStdin streams data from the local process to the remote process.
	streamStdin(conn *GCCon)
}

// GCArguments represents a set of arguments for a GoneCat connection.
type GCArguments struct {
	// AddrStr is the address string of the remote process.
	AddrStr string
	// AddrPort is the port number of the remote process.
	AddrPort string
	// Network specifies the network protocol to use (e.g., "tcp", "udp").
	Network string
	// Listening indicates whether the connection should be in listening mode.
	Listening bool
	// Protocol specifies the protocol to use for the connection (e.g., "http", "ftp").
	Protocol string
	// IPVersion specifies the IP version to use (e.g., "4", "6").
	IPVersion string
	// SendCRLF indicates whether to send carriage return and line feed characters.
	SendCRLF bool
	// ReadStdin indicates whether to read from the standard input.
	ReadStdin bool
	// ReadPipe indicates whether to read from a pipe.
	ReadPipe bool
	// BufferSize specifies the size of the buffer to use for reading and writing.
	BufferSize int
	// Output is the writer to use for output.
	Output io.WriteCloser
}

func (gc *GCArguments) UseDefaults() {
	gc.Listening = false
	gc.IPVersion = ""
	gc.Protocol = "tcp"
	gc.SendCRLF = true
	gc.ReadStdin = true
	gc.ReadPipe = false
	gc.BufferSize = 1024
	gc.Output = os.Stdout
}

// GetCat returns a GoneCat based on the given arguments.
func GetCat(gc *GCArguments) GoneCat {
	addr := gc.resolveAddress()

	switch gc.Protocol {
	case "tcp":
		return &TcpCat{gc, addr.(*net.TCPAddr)}
	case "udp":
		return &UdpCat{gc, addr.(*net.UDPAddr)}
	default:
		log.Fatalf("Wrong protocol name %s", gc.Protocol)
		return nil
	}
}

// resolveAddress resolves the given address and protocol into a net.Addr.
func (gc *GCArguments) resolveAddress() net.Addr {
	gc.Network = gc.Protocol + gc.IPVersion

	ip := net.ParseIP(gc.AddrStr)
	port, err := strconv.Atoi(gc.AddrPort)
	if err != nil {
		log.Fatal("The given port is not a number")
	}

	if gc.Protocol == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port, Zone: ""}
	} else if gc.Protocol == "udp" {
		return &net.UDPAddr{IP: ip, Port: port, Zone: ""}
	} else {
		return nil
	}
}
