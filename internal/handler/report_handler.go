package handler

import (
	"encoding/json"
	"net/http"
	"netops-api/internal/domain"
	"netops-api/internal/logging"
	"netops-api/internal/repository"
)

type ReportHandler struct {
	logger *logging.Logger
	repo   *repository.ReportRepository
}

func NewReportHandler(logger *logging.Logger, repo *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{
		logger: logger,
		repo:   repo,
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

	if err := h.repo.SaveReport(r.Context(), req); err != nil {
		h.logger.Error("échec de la sauvegarde en DB pour la machine (%s): %v", req.DeviceID, err)
		respondWithError(w, http.StatusServiceUnavailable, "erreur lors de l'insertion en base")
		return
	}

	h.logger.Info("nouveau rapport reçu pour la machine (%s)", req.DeviceID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	reports, err := h.repo.GetReports(r.Context())
	if err != nil {
		h.logger.Error("échec de la récupération des historiques: %v", err)
		respondWithError(w, http.StatusServiceUnavailable, "erreur lors de la récupération des historiques")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := h.repo.GetReportByID(r.Context(), id)
	if err != nil {
		h.logger.Error("échec de la récupération de l'historique (%s): %v", id, err)
		respondWithError(w, http.StatusServiceUnavailable, "erreur lors de la récupération du report "+id)
		return
	}

	if report == nil {
		respondWithError(w, http.StatusNotFound, "report introuvable")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(report)
}
