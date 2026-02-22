package scanner

import (
	"net"
	"time"
)

// Delegate to allow injecting mocks for testing purposes.
type DialerFunc func(network, address string, timeout time.Duration) (net.Conn, error)

// Performs TCP reachability checks on hosts.
type TCPPortChecker struct {
	dial DialerFunc
}

// NewTCPPortChecker creates a scanner using the real network dialer.
func NewTCPPortChecker() *TCPPortChecker {
	return &TCPPortChecker{
		dial: net.DialTimeout,
	}
}

// NewTCPPortCheckerWithDial allows tests injecting a custom dial function.
func NewTCPPortCheckerWithDial(dial DialerFunc) *TCPPortChecker {
	return &TCPPortChecker{
		dial: dial,
	}
}

// PingIP checks if a host is reachable by attempting a TCP connection.
func (n *TCPPortChecker) TCPPing(ip string, timeout time.Duration) bool {
	conn, err := n.dial("tcp", ip+":80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
