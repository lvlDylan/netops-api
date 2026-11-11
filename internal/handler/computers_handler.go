package handler

import (
	"encoding/json"
	"net/http"
	"netops-api/internal/domain"
	"netops-api/internal/logging"
)

type ComputerHandler struct {
	logger *logging.Logger
}

func NewComputerHandler(logger *logging.Logger) *ComputerHandler {
	return &ComputerHandler{
		logger: logger,
	}
}

func (h *ComputerHandler) ListComputers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode([]map[string]string{
		{"id": "6b1a3ec7-67cb-4271-aa3e-db43488b18c1", "hostname": "windows-dylan"},
	})
	if err != nil {

	}
}

func (h *ComputerHandler) GetComputerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"id":       id,
		"hostname": "windows-dylan",
	})
	if err != nil {

	}
}

func (h *ComputerHandler) InsertReport(w http.ResponseWriter, r *http.Request) {
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
