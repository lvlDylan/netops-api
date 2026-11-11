package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"netops-api/internal/logging"
)

type HealthHandler struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewHealthHandler(logger *logging.Logger, db *sql.DB) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		db:     db,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {

	if err := h.db.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ko",
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
