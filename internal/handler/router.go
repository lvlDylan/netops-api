package handler

import (
	"encoding/json"
	"net/http"
	"netops-api/internal/logging"
	"netops-api/internal/middleware"
)

func NewRouter(logger *logging.Logger,
	healthHandler *HealthHandler, computerHandler *ComputerHandler,
	reportHandler *ReportHandler,
) http.Handler {
	mainMux := http.NewServeMux()

	mainMux.HandleFunc("GET /health", healthHandler.Check)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /computers", computerHandler.ListComputers)
	apiMux.HandleFunc("GET /computers/{id}", computerHandler.GetComputerByID)

	apiMux.HandleFunc("GET /reports", reportHandler.GetReport)
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
