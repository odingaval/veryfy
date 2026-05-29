package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/odingaval/veryfy/api/internal/config"
	"github.com/odingaval/veryfy/api/internal/handlers"
	"github.com/odingaval/veryfy/api/internal/middleware"
	"github.com/odingaval/veryfy/api/internal/repositories"
	"github.com/odingaval/veryfy/api/internal/services"
)

type repositorySet interface {
	repositories.LicenseRepository
	repositories.IssuerRepository
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	var repo repositorySet = repositories.NewMemoryRepository()
	if cfg.SolanaRPCURL != "" && cfg.ProgramID != "" {
		solanaRepo, err := repositories.NewSolanaRepository(cfg)
		if err != nil {
			logger.Error("failed to initialize solana repository", "error", err)
			os.Exit(1)
		}
		repo = solanaRepo
		logger.Info("using solana repository", "cluster", cfg.SolanaCluster, "program_id", cfg.ProgramID)
	} else {
		logger.Info("using in-memory repository")
	}

	hashService := services.NewHashService()
	licenseService := services.NewLicenseService(hashService, repo)
	issuerService := services.NewIssuerService(repo)

	router := http.NewServeMux()
	router.Handle("/health", handlers.NewHealthHandler(time.Now()))
	router.Handle("/issuers/register", handlers.NewRegisterIssuerHandler(issuerService))
	router.Handle("/licenses/issue", handlers.NewIssueLicenseHandler(licenseService))
	router.Handle("/licenses/verify", handlers.NewVerifyLicenseHandler(licenseService))
	router.Handle("/licenses/revoke", handlers.NewRevokeLicenseHandler(licenseService))

	handler := middleware.Recover(logger)(
		middleware.RequestLogger(logger)(
			middleware.CORS(cfg.AllowedOrigins)(
				middleware.Timeout(cfg.RequestTimeout)(router),
			),
		),
	)

	server := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
}
