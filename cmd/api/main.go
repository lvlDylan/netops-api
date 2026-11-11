package main

import (
	"log"
	"net/http"
	"netops-api/internal/config"
	"netops-api/internal/handler"
	"netops-api/internal/logging"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("erreur config: %v", err)
	}

	logger, err := logging.NewLogger(cfg.Verbose, cfg.LogFile)
	if err != nil {
		log.Fatalf("impossible d'initialiser le logger")
	}
	defer logger.Close()

	logger.Info("démarrage du serveur api")

	db, err := config.NewDatabase(cfg.DSN())
	if err != nil {
		logger.Error("échec connexion base de données %w", err)
		log.Fatalf("impossible de se connecter à la base de données")
	}
	defer db.Close()
	logger.Debug("connexion bd ok")

	computerHandler := handler.NewComputerHandler(logger)
	healthHandler := handler.NewHealthHandler(logger, db)

	router := handler.NewRouter(logger, computerHandler, healthHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("impossible d'écouter sur le port 8080: %v", err)
		return
	}
}
