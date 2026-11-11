package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"netops-api/internal/logging"
	"netops-api/internal/middleware"
)

func main() {
	verbose := flag.Bool("v", false, "Activer les logs debugs")
	logFile := flag.Bool("log-file", true, "Activer la sortie de log sur le fichier log temporaire")
	flag.Parse()

	logger, err := logging.NewLogger(*verbose, *logFile)
	if err != nil {
		log.Fatalf("impossible d'initialiser le logger")
	}
	defer logger.Close()

	logger.Info("démarrage du serveur api")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "OK",
			Message: "Hello World",
		})
	})

	logMiddleware := middleware.LoggerMiddleware(logger)
	wrappedMux := logMiddleware(mux)

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Errorf("impossible d'écouter sur le port 8080")
		return
	}
}
