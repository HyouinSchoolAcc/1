package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"data_labler_ui_go/internal/web"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}

	cfg := web.Config{
		RootDir:         root,
		TemplatesDir:    filepath.Join(root, "templates"),
		StaticDir:       filepath.Join(root, "static"),
		PresetBaseDir:   filepath.Join(root, "presets"), // Use presets from project directory
		BackupRootDir:   filepath.Join(root, "backups"),
		RetentionDays:   7,
		ServePort:       getenvDefault("SERVE_PORT", "5003"),
		EnableDevReload: true,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	web.RegisterRoutes(r, cfg)

	addr := "0.0.0.0:" + cfg.ServePort
	log.Printf("Go server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
