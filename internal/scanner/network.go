package scanner

import (
	"net"
	"time"
)

/******** TCP Ping ********/

type DialerFunc func(network, address string, timeout time.Duration) (net.Conn, error)

type TCPPortChecker struct {
	dial DialerFunc
}

func NewTCPPortChecker() *TCPPortChecker {
	return &TCPPortChecker{
		dial: net.DialTimeout,
	}
}

func NewTCPPortCheckerWithDial(dial DialerFunc) *TCPPortChecker {
	return &TCPPortChecker{
		dial: dial,
	}
}

func (n *TCPPortChecker) Ping(ip string, timeout time.Duration) bool {
	conn, err := n.dial("tcp", ip+":80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
