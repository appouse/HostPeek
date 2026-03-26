package collector

import (
	"github.com/shirou/gopsutil/v4/mem"
)

// MemoryInfo contains memory metrics.
type MemoryInfo struct {
	TotalMB      uint64  `json:"total_mb"`
	UsedMB       uint64  `json:"used_mb"`
	FreeMB       uint64  `json:"free_mb"`
	UsagePercent float64 `json:"usage_percent"`
}

// CollectMemory gathers memory metrics.
func CollectMemory() (*MemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		TotalMB:      v.Total / 1024 / 1024,
		UsedMB:       v.Used / 1024 / 1024,
		FreeMB:       v.Available / 1024 / 1024,
		UsagePercent: round2(v.UsedPercent),
	}, nil
}
