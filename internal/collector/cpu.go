package collector

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
)

// CPUInfo contains CPU metrics.
type CPUInfo struct {
	Cores        int       `json:"cores"`
	ModelName    string    `json:"model_name"`
	UsagePercent float64   `json:"usage_percent"`
	LoadAvg      []float64 `json:"load_avg"`
}

// CollectCPU gathers CPU metrics.
func CollectCPU() (*CPUInfo, error) {
	cores, err := cpu.Counts(true) // logical cores
	if err != nil {
		return nil, err
	}

	percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return nil, err
	}

	var usagePercent float64
	if len(percentages) > 0 {
		usagePercent = round2(percentages[0])
	}

	// Get CPU model name
	var modelName string
	infos, err := cpu.Info()
	if err == nil && len(infos) > 0 {
		modelName = infos[0].ModelName
	}

	// Load averages (Linux/macOS only; returns zeros on Windows)
	loadAvg := []float64{0, 0, 0}
	if runtime.GOOS != "windows" {
		avg, err := load.Avg()
		if err == nil && avg != nil {
			loadAvg = []float64{
				round2(avg.Load1),
				round2(avg.Load5),
				round2(avg.Load15),
			}
		}
	}

	return &CPUInfo{
		Cores:        cores,
		ModelName:    modelName,
		UsagePercent: usagePercent,
		LoadAvg:      loadAvg,
	}, nil
}

func round2(v float64) float64 {
	return float64(int(v*100)) / 100
}
