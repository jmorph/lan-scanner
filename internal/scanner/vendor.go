package scanner

import (
	"bufio"
	"io"
	"strings"
)

type VendorDB struct {
	data map[string]string
}

// creator of dictionary of vendors
func NewVendorDB() *VendorDB {
	return &VendorDB{
		data: make(map[string]string),
	}
}

// load vendors from text source
func (v *VendorDB) LoadFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty line
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		prefix := normalizePrefix(parts[0])
		v.data[prefix] = strings.TrimSpace(parts[1])
	}

	return scanner.Err()
}

// find vendor by mac
func (v *VendorDB) Lookup(mac string) string {
	prefix := extractPrefix(mac)
	if vendor, ok := v.data[prefix]; ok {
		return vendor
	}
	return "Unknown"
}

func normalizePrefix(p string) string {
	return strings.ToUpper(strings.ReplaceAll(p, ":", ""))
}

// cleans mac and takes first 6 hex characters (OUI - Organizationally Unique Identifier)
func extractPrefix(mac string) string {
	clean := strings.ToUpper(strings.ReplaceAll(mac, ":", ""))
	if len(clean) < 6 {
		return ""
	}
	return clean[:6]
}
