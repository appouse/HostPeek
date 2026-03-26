package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

// UptimeInfo contains system uptime information.
type UptimeInfo struct {
	UptimeSeconds uint64 `json:"uptime_seconds"`
	BootTime      string `json:"boot_time"`
}

// CollectUptime gathers uptime information.
func CollectUptime() (*UptimeInfo, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return nil, err
	}

	bootTime, err := host.BootTime()
	if err != nil {
		return nil, err
	}

	return &UptimeInfo{
		UptimeSeconds: uptime,
		BootTime:      time.Unix(int64(bootTime), 0).UTC().Format(time.RFC3339),
	}, nil
}
