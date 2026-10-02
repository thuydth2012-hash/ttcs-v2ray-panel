package util

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemInfo struct {
	CPUPercent      float64 `json:"cpu_percent"`
	MemoryPercent   float64 `json:"memory_percent"`
	DiskPercent     float64 `json:"disk_percent"`
	SwapPercent     float64 `json:"swap_percent"`
	Uptime          uint64  `json:"uptime"`
	LoadAverage     float64 `json:"load_average"`
	UploadTraffic   uint64  `json:"upload_traffic"`
	DownloadTraffic uint64  `json:"download_traffic"`
	UploadSpeed     float64 `json:"upload_speed"`
	DownloadSpeed   float64 `json:"download_speed"`
	ActiveProcesses int     `json:"active_processes"`
}

var (
	lastUploadBytes   uint64
	lastDownloadBytes uint64
	lastMeasureTime   time.Time
)

func GetSystemInfo() (*SystemInfo, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, err
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	activeProcesses := len(processes)
	interfaces, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var uploadTraffic uint64
	var downloadTraffic uint64

	for _, iface := range interfaces {
		if iface.Name == "ens33" {
			uploadTraffic = iface.BytesSent
			downloadTraffic = iface.BytesRecv
			break
		}
	}

	now := time.Now()

	var uploadSpeed float64
	var downloadSpeed float64

	if !lastMeasureTime.IsZero() {
		elapsed := now.Sub(lastMeasureTime).Seconds()

		if elapsed > 0 && uploadTraffic >= lastUploadBytes && downloadTraffic >= lastDownloadBytes {
			uploadSpeed = float64(uploadTraffic-lastUploadBytes) / elapsed
			downloadSpeed = float64(downloadTraffic-lastDownloadBytes) / elapsed
		}
	}

	lastUploadBytes = uploadTraffic
	lastDownloadBytes = downloadTraffic
	lastMeasureTime = now

	cpuUsage := float64(0)

	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	return &SystemInfo{
		CPUPercent:      cpuUsage,
		MemoryPercent:   memory.UsedPercent,
		DiskPercent:     diskInfo.UsedPercent,
		SwapPercent:     swap.UsedPercent,
		Uptime:          uptime,
		LoadAverage:     loadAvg.Load1,
		UploadTraffic:   uploadTraffic,
		DownloadTraffic: downloadTraffic,
		UploadSpeed:     uploadSpeed,
		DownloadSpeed:   downloadSpeed,
		ActiveProcesses: activeProcesses,
	}, nil
}
