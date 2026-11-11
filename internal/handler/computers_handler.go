package handler

import (
	"encoding/json"
	"net/http"
)

type ComputerHandler struct {
}

func NewComputerHandler() *ComputerHandler {
	return &ComputerHandler{}
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
