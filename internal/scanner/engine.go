package scanner

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

/******** Interfaces ********/

type Pinger interface {
	Ping(ip string, timeout time.Duration) bool
}

type ARPProvider interface {
	GetARP() map[string]string
}

type VendorLookup interface {
	Lookup(mac string) string
}

/******** Engine ********/

type Engine struct {
	pinger  Pinger
	arp     ARPProvider
	vendor  VendorLookup
	workers int
}

func NewEngine(p Pinger, a ARPProvider, v VendorLookup, workers int) *Engine {
	if workers <= 0 {
		workers = 50
	}
	return &Engine{
		pinger:  p,
		arp:     a,
		vendor:  v,
		workers: workers,
	}
}

func (e *Engine) ScanSubnet(
	ctx context.Context,
	baseIP string,
	timeout time.Duration,
) ScanResult {

	jobs := make(chan string)
	results := make(chan Device)

	var wg sync.WaitGroup

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if e.pinger.Ping(ip, timeout) {
					results <- Device{IP: ip}
				}
			}
		}()
	}

	go func() {
		for i := 1; i <= 254; i++ {
			jobs <- baseIP + "." + strconv.Itoa(i)
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var devices []Device
	for d := range results {
		devices = append(devices, d)
	}

	e.enrich(devices)

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].IP < devices[j].IP
	})

	return ScanResult{
		Devices:  devices,
		ScanTime: time.Now().UTC().Format(time.RFC3339),
	}
}

func (e *Engine) enrich(devices []Device) {
	arp := e.arp.GetARP()

	for i := range devices {
		if mac, ok := arp[devices[i].IP]; ok {
			devices[i].MAC = mac
			devices[i].Vendor = e.vendor.Lookup(mac)

			if strings.Contains(strings.ToLower(devices[i].Vendor), "camera") {
				devices[i].DeviceType = DeviceIoT
			} else {
				devices[i].DeviceType = DeviceUnknown
			}
		}
	}
}
