package domain

import "time"

type Report struct {
	ID              int       `json:"id"`
	SendAt          time.Time `json:"send_at"`
	DeviceID        string    `json:"device_id"`
	IPAddress       string    `json:"ip_address"`
	AgentVersion    string    `json:"agent_version"`
	UptimeSecond    int       `json:"uptime_second"`
	MemUsagePercent float64   `json:"mem_usage_percent"`
	MemUsedBytes    uint64    `json:"mem_used_bytes"`
	MemTotalBytes   uint64    `json:"mem_total_bytes"`
	CPUUsagePercent float64   `json:"cpu_usage_percent"`
	CPUCores        int       `json:"cpu_cores"`
}
