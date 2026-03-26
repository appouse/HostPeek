package collector

import (
	"net"

	psnet "github.com/shirou/gopsutil/v4/net"
)

// NetworkInfo contains network metrics.
type NetworkInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
}

// InterfaceInfo contains information about a network interface.
type InterfaceInfo struct {
	Name    string   `json:"name"`
	IPv4    []string `json:"ipv4"`
	IPv6    []string `json:"ipv6,omitempty"`
	MAC     string   `json:"mac"`
	RxBytes uint64   `json:"rx_bytes"`
	TxBytes uint64   `json:"tx_bytes"`
}

// CollectNetwork gathers network interface metrics.
func CollectNetwork() (*NetworkInfo, error) {
	ifaces, err := psnet.Interfaces()
	if err != nil {
		return nil, err
	}

	// Get I/O counters per interface
	counters, err := psnet.IOCounters(true)
	if err != nil {
		return nil, err
	}

	counterMap := make(map[string]psnet.IOCountersStat)
	for _, c := range counters {
		counterMap[c.Name] = c
	}

	var interfaces []InterfaceInfo
	for _, iface := range ifaces {
		// Skip loopback and interfaces without addresses
		if iface.Flags != nil {
			isLoopback := false
			for _, f := range iface.Flags {
				if f == "loopback" {
					isLoopback = true
					break
				}
			}
			if isLoopback {
				continue
			}
		}

		var ipv4, ipv6 []string
		for _, addr := range iface.Addrs {
			ip, _, err := net.ParseCIDR(addr.Addr)
			if err != nil {
				continue
			}
			if ip.To4() != nil {
				ipv4 = append(ipv4, ip.String())
			} else {
				ipv6 = append(ipv6, ip.String())
			}
		}

		// Skip interfaces without any IP
		if len(ipv4) == 0 && len(ipv6) == 0 {
			continue
		}

		info := InterfaceInfo{
			Name: iface.Name,
			IPv4: ipv4,
			IPv6: ipv6,
			MAC:  iface.HardwareAddr,
		}

		if c, ok := counterMap[iface.Name]; ok {
			info.RxBytes = c.BytesRecv
			info.TxBytes = c.BytesSent
		}

		interfaces = append(interfaces, info)
	}

	return &NetworkInfo{Interfaces: interfaces}, nil
}
