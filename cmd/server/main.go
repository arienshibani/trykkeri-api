package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trykkeri-api/internal/config"
	"trykkeri-api/internal/handler"
	"trykkeri-api/internal/middleware"
	"trykkeri-api/internal/pdf"
)

const version = "1.0.0"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	initLogging(cfg.JSONLogs)

	startTime := time.Now()
	pdfSvc := pdf.NewService(cfg)
	h := handler.New(cfg, pdfSvc, version, startTime)

	router := handler.Routes(h)
	wrapped := middleware.Chain(router, cfg, version)

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: wrapped}

	go func() {
		slog.Info("Starting HTML→PDF API server", "version", version, "port", cfg.Port)
		slog.Info("Server listening", "address", addr)
		slog.Info("Docs", "url", "http://localhost:"+fmt.Sprint(cfg.Port)+"/openapi.json")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Received shutdown signal, shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "err", err)
	}
	slog.Info("Server shut down gracefully")
}

func initLogging(jsonLogs bool) {
	var handler slog.Handler
	if jsonLogs {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	slog.SetDefault(slog.New(handler))
}
