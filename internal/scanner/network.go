package scanner

import (
	"net"
	"time"
)

// PingIP performs a TCP connect to port 80 for fast cross-platform ping
func PingIP(ip string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", ip+":80", timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
