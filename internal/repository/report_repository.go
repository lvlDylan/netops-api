package repository

import (
	"context"
	"database/sql"
	"fmt"
	"netops-api/internal/domain"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReport(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) SaveReport(ctx context.Context, req domain.ComputerRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("impossible de démarrer la transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	deviceQuery := `
			INSERT INTO device (id, hostname, os, platform, arch)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
                     hostname = VALUES(hostname),
                     os = VALUES(os),
                     platform = VALUES(platform),
                     arch = VALUES(arch)
	`

	_, err = tx.ExecContext(ctx, deviceQuery, req.DeviceID, req.Identity.Hostname, req.Identity.OS, req.Identity.Platform, req.Identity.Arch)
	if err != nil {
		return fmt.Errorf("échec de l'UPSERT device: %w", err)
	}

	historyQuery := `
		INSERT INTO history (
			send_at, device_id, ip_address, agent_version, uptime_seconds,
			mem_usage_percent, mem_used_bytes, mem_total_bytes,
			cpu_usage_percent, cpu_cores
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	res, err := tx.ExecContext(ctx, historyQuery,
		req.SendAt,
		req.DeviceID,
		req.Identity.IPAddress,
		req.Identity.AgentVersion,
		req.Identity.UptimeSeconds,
		req.Metrics.MemoryMetrics.UsagePercent,
		req.Metrics.MemoryMetrics.UsedBytes,
		req.Metrics.MemoryMetrics.TotalBytes,
		req.Metrics.CPUMetrics.UsagePercent,
		req.Metrics.CPUMetrics.Cores,
	)
	if err != nil {
		return fmt.Errorf("échec de l'INSERT history: %w", err)
	}

	historyID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("impossible de récupérer history_id: %w", err)
	}

	diskQuery := `
		INSERT INTO disk (history_id, mountpoint, usage_percent, used_bytes, total_bytes)
		VALUES (?, ?, ?, ?, ?);
	`
	stmt, err := tx.PrepareContext(ctx, diskQuery)
	if err != nil {
		return fmt.Errorf("erreur de préparation statement disk: %w", err)
	}
	defer stmt.Close()

	for _, d := range req.Metrics.DiskMetrics {
		_, err := stmt.ExecContext(ctx,
			historyID,
			d.MountPoint,
			d.UsagePercent,
			d.UsedBytes,
			d.TotalBytes,
		)
		if err != nil {
			return fmt.Errorf("échec d'insertion du disque %s: %w", d.MountPoint, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("échec du commit transaction: %w", err)
	}

	return nil
}
