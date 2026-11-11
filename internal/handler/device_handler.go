package handler

import (
	"encoding/json"
	"net/http"
	"netops-api/internal/logging"
	"netops-api/internal/repository"
)

type DeviceHandler struct {
	logger *logging.Logger
	repo   *repository.DeviceRepository
}

func NewDeviceHandler(logger *logging.Logger, repo *repository.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.repo.GetDevices(r.Context())
	if err != nil {
		h.logger.Error("échec de la récupération des devices: %v", err)
		respondWithError(w, http.StatusServiceUnavailable, "erreur lors de la récupération des machines")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(devices)
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	device, err := h.repo.GetDeviceByID(r.Context(), id)
	if err != nil {
		h.logger.Error("échec de la récupération du device (%s): %v", id, err)
		respondWithError(w, http.StatusServiceUnavailable, "erreur lors de la récupération de la machine "+id)
		return
	}

	if device == nil {
		respondWithError(w, http.StatusNotFound, "machine introuvable")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(device)
}
