// Package logging fournit un logger simple avec gestion de niveaux (INFO, ERROR, DEBUG)
// et support de la sortie fichier.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger gère l'écriture des logs vers stdout/stderr (et un fichier).
type Logger struct {
	verbose  bool
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
	logFile  *os.File
}

// NewLogger instancie un nouveau Logger.
// Si writeToFile est true, les logs sont aussi écrits dans le dossier temp de l'OS.
func NewLogger(verbose bool, writeToFile bool) (*Logger, error) {
	flags := log.Ldate | log.Ltime | log.Lmsgprefix

	var writer io.Writer = os.Stdout
	var errWriter io.Writer = os.Stderr
	var logFile *os.File

	if writeToFile {
		logPath := filepath.Join(os.TempDir(), "netops", "netops-agent.log")
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("erreur lors de la creation du repertoire %s: %w", logDir, err)
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la creation du fichier de log %w", err)
		}

		logFile = file
		writer = io.MultiWriter(os.Stdout, logFile)
		errWriter = io.MultiWriter(os.Stderr, logFile)
	}

	return &Logger{
		verbose:  verbose,
		infoLog:  log.New(writer, "[INFO] ", flags),
		errorLog: log.New(errWriter, "[ERROR] ", flags),
		debugLog: log.New(writer, "[DEBUG] ", flags),
		logFile:  logFile,
	}, nil
}

// Info écrit un log d'information standard.
func (l *Logger) Info(format string, args ...any) {
	l.infoLog.Printf(format, args...)
}

// Error écrit un log d'erreur sur stderr.
func (l *Logger) Error(format string, args ...any) {
	l.errorLog.Printf(format, args...)
}

// Debug écrit un log de debug uniquement si le mode verbose est actif.
func (l *Logger) Debug(format string, args ...any) {
	if l.verbose {
		l.debugLog.Printf(format, args...)
	}
}

// Close ferme le fichier de log s'il est ouvert.
func (l *Logger) Close() {
	if l.logFile != nil {
		_ = l.logFile.Close()
	}
}
