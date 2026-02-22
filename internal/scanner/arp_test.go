package scanner

import (
	"io"
	"strings"
	"testing"
)

/******* Parser *******/

// Test ARP output parser
func TestParseARP(t *testing.T) {
	mockARP := `
				192.168.1.1    ether   aa:bb:cc:dd:ee:ff
				192.168.1.10   ether   11:22:33:44:55:66
				`

	result := parseARP(strings.NewReader(mockARP))

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}

	if result["192.168.1.1"] != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("wrong MAC for 192.168.1.1")
	}

	if result["192.168.1.10"] != "11:22:33:44:55:66" {
		t.Errorf("wrong MAC for 192.168.1.10")
	}
}

/******* Command Runner *******/

type mockRunner struct {
	output string
	err    error
}

func (m mockRunner) Run(name string, args ...string) (io.ReadCloser, error) {
	if m.err != nil {
		return nil, m.err
	}
	return io.NopCloser(strings.NewReader(m.output)), nil
}

func TestGetARPTable(t *testing.T) {
	mockOutput := `
					192.168.1.10    ether   ff:ee:dd:cc:bb:aa
					`

	scanner := &ARPScanner{
		runner: mockRunner{output: mockOutput},
	}

	result, err := scanner.GetARPTable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["192.168.1.10"] != "ff:ee:dd:cc:bb:aa" {
		t.Errorf("incorrect parsing from GetARPTable")
	}
}
