package scanner

type DeviceType string

const (
	DeviceUnknown DeviceType = "Unknown"
	DeviceIoT     DeviceType = "IoT"
	DevicePC      DeviceType = "PC"
	DevicePhone   DeviceType = "Phone"
	DeviceNAS     DeviceType = "NAS"
)

type Device struct {
	IP         string     `json:"ip"`
	MAC        string     `json:"mac"`
	Vendor     string     `json:"vendor"`
	DeviceType DeviceType `json:"deviceType"`
}

type ScanResult struct {
	Devices  []Device `json:"devices"`
	ScanTime string   `json:"scanTime"`
}
