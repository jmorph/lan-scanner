package scanner

import (
	"strings"
	"testing"
)

func TestLoadAndLookupVendor(t *testing.T) {
	mockData := `
				AA11BB,VendorOne
				CC22DD,VendorTwo
				`

	db := NewVendorDB()

	err := db.LoadFromReader(strings.NewReader(mockData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// list of vendors for assert
	tests := []struct {
		mac      string
		expected string
	}{
		{"aa:11:bb:ff:ee:dd", "VendorOne"},
		{"CC:22:DD:00:11:22", "VendorTwo"},
		{"00:00:00:00:00:00", "Unknown"},
		{"invalid", "Unknown"},
	}

	for _, tt := range tests {
		result := db.Lookup(tt.mac)
		if result != tt.expected {
			t.Errorf("for MAC %s expected %s got %s",
				tt.mac, tt.expected, result)
		}
	}
}
