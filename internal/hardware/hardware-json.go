package hardware

import (
	"encoding/json"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemInfo struct {
	OperatingSystem       string             `json:"operating_system"`
	Platform              string             `json:"platform"`
	Hostname              string             `json:"hostname"`
	Processes             uint64             `json:"processes"`
	Memory                MemoryInfo         `json:"memory"`
	Disk                  DiskInfo           `json:"disk"`
	CPU                   CPUInfo            `json:"cpu"`
}

type MemoryInfo struct {
	TotalMB        uint64  `json:"total_mb"`
	FreeMB         uint64  `json:"free_mb"`
	UsedPercent    float64 `json:"used_percent"`
}

type DiskInfo struct {
	TotalGB       uint64  `json:"total_gb"`
	UsedGB        uint64  `json:"used_gb"`
	FreeGB        uint64  `json:"free_gb"`
	UsedPercent   float64 `json:"used_percent"`
}

type CPUInfo struct {
	ModelName     string    `json:"model_name"`
	Family        string    `json:"family"`
	SpeedMHz      float64   `json:"speed_mhz"`
	CoreUsage     []float64 `json:"core_usage"`
}

// GetSystemInfoJSON retrieves system information in JSON format
func GetSystemInfoJSON() (string, error) {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return "", err
	}

	cpuInfo := CPUInfo{}
	if len(cpuStat) != 0 {
		cpuInfo.ModelName = cpuStat[0].ModelName
		cpuInfo.Family = cpuStat[0].Family
		cpuInfo.SpeedMHz = cpuStat[0].Mhz
		cpuInfo.CoreUsage = percentage
	}

	systemInfo := SystemInfo{
		OperatingSystem: runTimeOS,
		Platform:        hostStat.Platform,
		Hostname:        hostStat.Hostname,
		Processes:       hostStat.Procs,
		Memory: MemoryInfo{
			TotalMB:     vmStat.Total / megabyteDiv,
			FreeMB:      vmStat.Free / megabyteDiv,
			UsedPercent: vmStat.UsedPercent,
		},
		Disk: DiskInfo{
			TotalGB:     diskStat.Total / gigabyteDiv,
			UsedGB:      diskStat.Used / gigabyteDiv,
			FreeGB:      diskStat.Free / gigabyteDiv,
			UsedPercent: diskStat.UsedPercent,
		},
		CPU: cpuInfo,
	}

	jsonData, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
