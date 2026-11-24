package domain

type DiskAverage struct {
	Mountpoint   string `json:"mountpoint"`
	UsagePercent string `json:"usage_percent"`
}

type Average struct {
	DeviceID        string        `json:"device_id"`
	PeriodStart     string        `json:"period_start"`
	MemUsagePercent string        `json:"os"`
	CpuUsagePercent string        `json:"platform"`
	Disks           []DiskAverage `json:"disks"`
}
