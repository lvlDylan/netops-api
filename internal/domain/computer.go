package domain

import "time"

type Identity struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Platform      string `json:"platform"`
	Arch          string `json:"arch"`
	IPAddress     string `json:"ip_address"`
	AgentVersion  string `json:"agent_version"`
	UptimeSeconds uint64 `json:"uptime_seconds"`
}

type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

type MemoryMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	UsedBytes    uint64  `json:"used_bytes"`
	TotalBytes   uint64  `json:"total_bytes"`
}

type DiskMetrics struct {
	MountPoint   string  `json:"mount_point"`
	UsagePercent float64 `json:"usage_percent"`
	UsedBytes    uint64  `json:"used_bytes"`
	TotalBytes   uint64  `json:"total_bytes"`
}

type Metrics struct {
	CPUMetrics    CPUMetrics    `json:"cpu"`
	MemoryMetrics MemoryMetrics `json:"memory"`
	DiskMetrics   []DiskMetrics `json:"disk"`
}

type ComputerRequest struct {
	DeviceID string    `json:"device_id"`
	SendAt   time.Time `json:"send_at"`
	Identity Identity  `json:"identity"`
	Metrics  Metrics   `json:"metrics"`
}
