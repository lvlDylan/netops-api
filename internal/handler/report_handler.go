package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"netops-api/internal/domain"
	"netops-api/internal/logging"
)

type ReportHandler struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewReportHandler(logger *logging.Logger, db *sql.DB) *ReportHandler {
	return &ReportHandler{
		logger: logger,
		db:     db,
	}
}

func (h *ReportHandler) InsertReport(w http.ResponseWriter, r *http.Request) {
	var req domain.ComputerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "payload invalide")
		return
	}

	if req.DeviceID == "" {
		respondWithError(w, http.StatusBadRequest, "payload invalide")
		return
	}

	h.logger.Info("nouveau rapport reçu pour la machine (%s)", req.DeviceID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}
