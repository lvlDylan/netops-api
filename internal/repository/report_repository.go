package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"netops-api/internal/domain"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) SaveReport(ctx context.Context, req domain.ReportRequest) error {
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

	reportQuery := `
		INSERT INTO report (
			send_at, device_id, ip_address, agent_version, uptime_seconds,
			mem_usage_percent, mem_used_bytes, mem_total_bytes,
			cpu_usage_percent, cpu_cores
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	res, err := tx.ExecContext(ctx, reportQuery,
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
		return fmt.Errorf("échec de l'INSERT report: %w", err)
	}

	reportID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("impossible de récupérer report_id: %w", err)
	}

	diskQuery := `
		INSERT INTO disk (report_id, mountpoint, usage_percent, used_bytes, total_bytes)
		VALUES (?, ?, ?, ?, ?);
	`
	stmt, err := tx.PrepareContext(ctx, diskQuery)
	if err != nil {
		return fmt.Errorf("erreur de préparation statement disk: %w", err)
	}
	defer stmt.Close()

	for _, d := range req.Metrics.DiskMetrics {
		_, err := stmt.ExecContext(ctx,
			reportID,
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

func (r *ReportRepository) GetReports(ctx context.Context) ([]domain.Report, error) {
	query := `SELECT id, send_at, device_id, ip_address, agent_version, uptime_seconds, mem_usage_percent, mem_used_bytes, mem_total_bytes, cpu_usage_percent, cpu_cores FROM report`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du SELECT report: %v", err)
	}
	defer rows.Close()

	var reports []domain.Report
	for rows.Next() {
		var report domain.Report
		err := rows.Scan(
			&report.ID, &report.SendAt, &report.DeviceID,
			&report.IPAddress, &report.AgentVersion, &report.UptimeSeconds,
			&report.MemUsagePercent, &report.MemUsedBytes, &report.MemTotalBytes,
			&report.CPUUsagePercent, &report.CPUCores)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan d'une ligne report: %w", err)
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération sur les lignes: %w", err)
	}

	return reports, nil
}

func (r *ReportRepository) GetReportByID(ctx context.Context, id string) (*domain.Report, error) {
	query := `SELECT id, send_at, device_id, ip_address, agent_version, uptime_seconds, mem_usage_percent, mem_used_bytes, mem_total_bytes, cpu_usage_percent, cpu_cores FROM report WHERE id = ?`
	var report domain.Report

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID, &report.SendAt, &report.DeviceID,
		&report.IPAddress, &report.AgentVersion, &report.UptimeSeconds,
		&report.MemUsagePercent, &report.MemUsedBytes, &report.MemTotalBytes,
		&report.CPUUsagePercent, &report.CPUCores,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur lors de la récupération d'un report: %w", err)
	}

	return &report, nil
}

func (r *ReportRepository) GetReportByDeviceID(ctx context.Context, deviceID string) ([]domain.Report, error) {
	query := `SELECT id, send_at, device_id, ip_address, agent_version, uptime_seconds, mem_usage_percent, mem_used_bytes, mem_total_bytes, cpu_usage_percent, cpu_cores FROM report WHERE device_id = ?`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du SELECT report: %v", err)
	}
	defer rows.Close()

	var reports []domain.Report
	for rows.Next() {
		var report domain.Report
		err := rows.Scan(
			&report.ID, &report.SendAt, &report.DeviceID,
			&report.IPAddress, &report.AgentVersion, &report.UptimeSeconds,
			&report.MemUsagePercent, &report.MemUsedBytes, &report.MemTotalBytes,
			&report.CPUUsagePercent, &report.CPUCores)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan d'une ligne report: %w", err)
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération sur les lignes: %w", err)
	}

	return reports, nil
}
