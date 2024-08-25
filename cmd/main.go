package main

import (
	"account-management/config"
	v1 "account-management/internal/handler/http/v1"
	"account-management/internal/repo"
	"account-management/internal/service"
	"account-management/pkg/client"
	"account-management/pkg/hasher"
	"account-management/pkg/logging"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func main() {
	logger := logging.GetLogger()

	cfg := config.NewConfig()

	db, err := client.NewClient(context.Background(), 3, cfg.DB)
	if err != nil {
		logger.Fatal("Error initializing db client: ", err)
	}

	repos := repo.NewRepositories(db)

	deps := service.ServiceDependencies{
		Repos:     repos,
		Hasher:    hasher.NewSHA256Hasher(),
		SecretKey: cfg.JWT.SecretKey,
		TokenTTL:  cfg.JWT.TokenTTL,
	}

	time.Sleep(5 * time.Second)

	services := service.NewServices(deps)

	router := chi.NewRouter()
	v1.NewHandler(router, services)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Infof("Starting server on %s", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed: ", err)
	}
}
