package scanner

import (
	"errors"
	"net"
	"testing"
	"time"
)

// dummy struck
type fakeConn struct{}

// implements all methods required by the net.Conn interface
func (f *fakeConn) Read(b []byte) (int, error)         { return 0, nil }
func (f *fakeConn) Write(b []byte) (int, error)        { return 0, nil }
func (f *fakeConn) Close() error                       { return nil }
func (f *fakeConn) LocalAddr() net.Addr                { return nil }
func (f *fakeConn) RemoteAddr() net.Addr               { return nil }
func (f *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (f *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (f *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func TestPingIP_Success(t *testing.T) {
	mockDial := func(network, address string, timeout time.Duration) (net.Conn, error) {
		return &fakeConn{}, nil
	}

	scanner := NewTCPPortCheckerWithDial(mockDial)

	result := scanner.Ping("192.168.1.100", 100*time.Millisecond)

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestPingIP_Failure(t *testing.T) {
	mockDial := func(network, address string, timeout time.Duration) (net.Conn, error) {
		return nil, errors.New("connection failed")
	}

	scanner := NewTCPPortCheckerWithDial(mockDial)

	result := scanner.Ping("192.168.1.100", 100*time.Millisecond)

	if result {
		t.Errorf("expected false, got true")
	}
}
