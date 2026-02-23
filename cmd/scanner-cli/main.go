package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmorph/lan-scanner/internal/scanner"
)

/******** Adapters ********/

type realPinger struct{}

func (r realPinger) Ping(ip string, timeout time.Duration) bool {
	return scanner.NewTCPPortChecker().Ping(ip, timeout)
}

type realARP struct{}

func (r realARP) GetARP() map[string]string {
	arp, _ := scanner.GetARPTable()
	return arp
}

type realVendor struct{}

func (r realVendor) Lookup(mac string) string {
	return scanner.LookupVendor(mac)
}

/******** Entry ********/

func main() {
	ctx := context.Background()

	engine := scanner.NewEngine(
		realPinger{},
		realARP{},
		realVendor{},
		50,
	)

	result := engine.ScanSubnet(ctx, "192.168.1", 500*time.Millisecond)

	for _, d := range result.Devices {
		fmt.Printf("%s | %s | %s | %s\n",
			d.IP,
			d.MAC,
			d.Vendor,
			d.DeviceType,
		)
	}
}
