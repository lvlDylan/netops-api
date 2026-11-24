package repository

import (
	"context"
	"database/sql"
	"fmt"
	"netops-api/internal/domain"
)

type AverageRepository struct {
	db *sql.DB
}

func (r *AverageRepository) GetAverageByDeviceID(ctx context.Context, deviceID string) (*domain.Average, error) {
	return nil, fmt.Errorf("not implemented yet")
}
