package lib

import "net"

// Returns a default local address, using localhost and a placeholder port
func defaultLocalAddr() net.Addr {
	// Use port 0 as a placeholder
	return &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 0,
	}
}
