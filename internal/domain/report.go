package domain

import "time"

type Report struct {
	ID              int       `json:"id"`
	SendAt          time.Time `json:"send_at"`
	DeviceID        string    `json:"device_id"`
	IPAddress       string    `json:"ip_address"`
	AgentVersion    string    `json:"agent_version"`
	UptimeSeconds   int       `json:"uptime_second"`
	MemUsagePercent float64   `json:"mem_usage_percent"`
	MemUsedBytes    uint64    `json:"mem_used_bytes"`
	MemTotalBytes   uint64    `json:"mem_total_bytes"`
	CPUUsagePercent float64   `json:"cpu_usage_percent"`
	CPUCores        int       `json:"cpu_cores"`
}

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

type ReportRequest struct {
	DeviceID string    `json:"device_id"`
	SendAt   time.Time `json:"send_at"`
	Identity Identity  `json:"identity"`
	Metrics  Metrics   `json:"metrics"`
}
