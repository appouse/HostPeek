package collector

import (
	"github.com/shirou/gopsutil/v4/host"
)

// OSInfo contains operating system information.
type OSInfo struct {
	Hostname     string `json:"hostname"`
	Platform     string `json:"platform"`
	Distribution string `json:"distribution"`
	Kernel       string `json:"kernel"`
}

// CollectOS gathers OS information.
func CollectOS() (*OSInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	distribution := info.Platform
	if info.PlatformVersion != "" {
		distribution += " " + info.PlatformVersion
	}

	return &OSInfo{
		Hostname:     info.Hostname,
		Platform:     info.OS,
		Distribution: distribution,
		Kernel:       info.KernelVersion,
	}, nil
}
