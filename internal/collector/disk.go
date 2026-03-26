package collector

import (
	"github.com/shirou/gopsutil/v4/disk"
)

// DiskInfo contains disk partition metrics.
type DiskInfo struct {
	Mount        string  `json:"mount"`
	FS           string  `json:"fs"`
	TotalGB      float64 `json:"total_gb"`
	UsedGB       float64 `json:"used_gb"`
	FreeGB       float64 `json:"free_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

// CollectDisk gathers disk partition metrics.
func CollectDisk() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // skip partitions we can't read
		}

		disks = append(disks, DiskInfo{
			Mount:        p.Mountpoint,
			FS:           p.Fstype,
			TotalGB:      round2(float64(usage.Total) / 1024 / 1024 / 1024),
			UsedGB:       round2(float64(usage.Used) / 1024 / 1024 / 1024),
			FreeGB:       round2(float64(usage.Free) / 1024 / 1024 / 1024),
			UsagePercent: round2(usage.UsedPercent),
		})
	}

	return disks, nil
}
