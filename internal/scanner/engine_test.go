package scanner

import (
	"context"
	"testing"
	"time"
)

/* ---------------- MOCKS ---------------- */

type mockPinger struct {
	live map[string]bool
}

func (m mockPinger) Ping(ip string, timeout time.Duration) bool {
	return m.live[ip]
}

type mockARP struct {
	data map[string]string
}

func (m mockARP) GetARP() map[string]string {
	return m.data
}

type mockVendor struct{}

func (m mockVendor) Lookup(mac string) string {
	if mac == "AA:BB:CC:DD:EE:FF" {
		return "CameraCorp"
	}
	return "GenericVendor"
}

/* ---------------- TESTS ---------------- */

func TestScanSubnet_BasicFlow(t *testing.T) {
	ctx := context.Background()

	pinger := mockPinger{
		live: map[string]bool{
			"192.168.1.10": true,
			"192.168.1.20": true,
		},
	}

	arp := mockARP{
		data: map[string]string{
			"192.168.1.10": "AA:BB:CC:DD:EE:FF",
			"192.168.1.20": "11:22:33:44:55:66",
		},
	}

	engine := NewEngine(pinger, arp, mockVendor{}, 5)

	result := engine.ScanSubnet(ctx, "192.168.1", time.Millisecond)

	if len(result.Devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(result.Devices))
	}

	// Ensure deterministic ordering
	if result.Devices[0].IP != "192.168.1.10" {
		t.Errorf("devices not sorted properly")
	}

	for _, d := range result.Devices {
		switch d.IP {
		case "192.168.1.10":
			if d.MAC != "AA:BB:CC:DD:EE:FF" {
				t.Errorf("wrong MAC for %s", d.IP)
			}
			if d.Vendor != "CameraCorp" {
				t.Errorf("wrong vendor for %s", d.IP)
			}
			if d.DeviceType != DeviceIoT {
				t.Errorf("expected IoT type for %s", d.IP)
			}

		case "192.168.1.20":
			if d.DeviceType != DeviceUnknown {
				t.Errorf("expected unknown type for %s", d.IP)
			}
		}
	}
}

func TestScanSubnet_NoARPEntry(t *testing.T) {
	ctx := context.Background()

	pinger := mockPinger{
		live: map[string]bool{
			"192.168.1.50": true,
		},
	}

	arp := mockARP{
		data: map[string]string{}, // no ARP entries
	}

	engine := NewEngine(pinger, arp, mockVendor{}, 5)

	result := engine.ScanSubnet(ctx, "192.168.1", time.Millisecond)

	if len(result.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(result.Devices))
	}

	d := result.Devices[0]

	if d.MAC != "" {
		t.Errorf("expected empty MAC")
	}
	if d.Vendor != "" {
		t.Errorf("expected empty Vendor")
	}
	if d.DeviceType != "" && d.DeviceType != DeviceUnknown {
		t.Errorf("unexpected device type")
	}
}

func TestScanSubnet_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	pinger := mockPinger{
		live: map[string]bool{
			"192.168.1.10": true,
		},
	}

	engine := NewEngine(pinger, mockARP{}, mockVendor{}, 5)

	result := engine.ScanSubnet(ctx, "192.168.1", time.Millisecond)

	if len(result.Devices) != 0 {
		t.Fatalf("expected 0 devices due to cancellation")
	}
}
