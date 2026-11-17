package handler

import (
	"encoding/json"
	"net/http"
	"netops-api/internal/logging"
	"netops-api/internal/middleware"
)

func NewRouter(logger *logging.Logger,
	healthHandler *HealthHandler, deviceHandler *DeviceHandler,
	reportHandler *ReportHandler,
) http.Handler {
	mainMux := http.NewServeMux()

	mainMux.HandleFunc("GET /health", healthHandler.Check)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /devices", deviceHandler.GetDevices)
	apiMux.HandleFunc("GET /devices/{id}", deviceHandler.GetDeviceByID)

	apiMux.HandleFunc("GET /reports", reportHandler.GetReport)
	apiMux.HandleFunc("GET /report/{id}", reportHandler.GetReportByID)
	apiMux.HandleFunc("GET /reports/{deviceID}", reportHandler.GetReportByDeviceID)
	apiMux.HandleFunc("POST /report", reportHandler.InsertReport)

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Not found",
			"path":  r.URL.Path,
		})
	})

	return middleware.LoggerMiddleware(logger)(mainMux)
}
