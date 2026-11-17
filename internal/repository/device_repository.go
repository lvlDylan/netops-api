package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"netops-api/internal/domain"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) GetDevices(ctx context.Context) ([]domain.Device, error) {
	query := `SELECT id, hostname, os, platform, arch FROM device`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du SELECT device: %v", err)
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var device domain.Device
		err := rows.Scan(&device.ID, &device.Hostname, &device.OS, &device.Platform, &device.Arch)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan d'une ligne device: %w", err)
		}

		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération sur les lignes: %w", err)
	}

	return devices, nil
}

func (r *DeviceRepository) GetOnlineDevices(ctx context.Context) ([]domain.Device, error) {
	query := `
    SELECT d.id, d.hostname, d.os, d.platform, d.arch
	FROM device d
	INNER JOIN (
		SELECT device_id, MAX(send_at) AS last_report FROM report GROUP BY device_id
	) latest
    ON latest.device_id = d.id
	WHERE latest.last_report > UTC_TIMESTAMP() - INTERVAL 4 MINUTE;
	`	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du SELECT device: %v", err)
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var device domain.Device
		err := rows.Scan(&device.ID, &device.Hostname, &device.OS, &device.Platform, &device.Arch)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan d'une ligne device: %w", err)
		}

		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération sur les lignes: %w", err)
	}

	return devices, nil
}

func (r *DeviceRepository) GetDeviceByID(ctx context.Context, id string) (*domain.Device, error) {
	query := `SELECT id, hostname, os, platform, arch FROM device WHERE id=?`
	var device domain.Device

	err := r.db.QueryRowContext(ctx, query, id).Scan(&device.ID, &device.Hostname, &device.OS, &device.Platform, &device.Arch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur lors de la récupération d'un device: %w", err)
	}

	return &device, nil
}
