package v1

import (
	"account-management/internal/service"
	"account-management/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func NewHandler(router chi.Router, services *service.Services) {
	router.Use(middleware.Logger)

	logger := logging.GetLogger()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Infof(`{"time":"%s", "method":"%s", "uri":"%s", "status":%d, "duration":"%s"}`, start.Format(time.RFC3339Nano), r.Method, r.RequestURI, http.StatusOK, duration)
		})
	})

	router.Use(middleware.Recoverer)

	// Health check route
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	auth := chi.NewRouter()
	{
		NewAuthRoutes(auth, services.Auth)
	}
	router.Mount("/auth", auth)
}
