package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// GetSystemStats gathers CPU and Memory usage.
func GetSystemStats() (map[string]interface{}, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	c, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get cpu stats: %w", err)
	}

	stats := map[string]interface{}{
		"cpu_percent":         c[0],
		"memory_total_mb":     v.Total / 1024 / 1024,
		"memory_used_mb":      v.Used / 1024 / 1024,
		"memory_used_percent": v.UsedPercent,
		"go_routines":         runtime.NumGoroutine(),
		"os":                  runtime.GOOS,
		"arch":                runtime.GOARCH,
	}

	return stats, nil
}

// GetProcessInfo gathers info about the current process.
func GetProcessInfo() (map[string]interface{}, error) {
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get process stats: %w", err)
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get process memory stats: %w", err)
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get process cpu stats: %w", err)
	}

	info := map[string]interface{}{
		"pid":           pid,
		"memory_rss_mb": memInfo.RSS / 1024 / 1024,
		"cpu_percent":   cpuPercent,
	}

	return info, nil
}
