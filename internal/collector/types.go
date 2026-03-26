package collector

import "github.com/Appouse/HostPeek/internal/config"

// FullMetrics is the combined response for GET /metrics.
type FullMetrics struct {
	Hostname string       `json:"hostname"`
	OS       *OSInfo      `json:"os,omitempty"`
	CPU      *CPUInfo     `json:"cpu,omitempty"`
	Memory   *MemoryInfo  `json:"memory,omitempty"`
	Disk     []DiskInfo   `json:"disk,omitempty"`
	Network  *NetworkInfo `json:"network,omitempty"`
	Uptime   *UptimeInfo  `json:"uptime,omitempty"`
	Agent    AgentInfo    `json:"agent"`
}

// AgentInfo contains agent metadata.
type AgentInfo struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// Version is set at build time via ldflags.
var Version = "dev"

// CollectAll gathers all enabled metrics into a single response.
func CollectAll(cfg *config.Config) (*FullMetrics, error) {
	m := &FullMetrics{
		Agent: AgentInfo{
			Version: Version,
			Status:  "ok",
		},
	}

	if cfg.Collectors.OS {
		osInfo, err := CollectOS()
		if err != nil {
			return nil, err
		}
		m.Hostname = osInfo.Hostname
		m.OS = osInfo
	}

	if cfg.Collectors.CPU {
		cpuInfo, err := CollectCPU()
		if err != nil {
			return nil, err
		}
		m.CPU = cpuInfo
	}

	if cfg.Collectors.Memory {
		memInfo, err := CollectMemory()
		if err != nil {
			return nil, err
		}
		m.Memory = memInfo
	}

	if cfg.Collectors.Disk {
		diskInfo, err := CollectDisk()
		if err != nil {
			return nil, err
		}
		m.Disk = diskInfo
	}

	if cfg.Collectors.Network {
		netInfo, err := CollectNetwork()
		if err != nil {
			return nil, err
		}
		m.Network = netInfo
	}

	if cfg.Collectors.Uptime {
		upInfo, err := CollectUptime()
		if err != nil {
			return nil, err
		}
		m.Uptime = upInfo
	}

	return m, nil
}
