package v1

import (
	"account-management/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AuthHandler struct {
	authService service.Auth
}

func NewAuthRoutes(router chi.Router, authService service.Auth) {
	authRoutes := &AuthHandler{
		authService: authService,
	}

	router.Route("/", func(r chi.Router) {
		r.Post("/sign-up", authRoutes.signUp)
		r.Post("/sign-in", authRoutes.signIn)
	})

}

type signUpInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var input signUpInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	userID, err := h.authService.CreateUser(r.Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
	})
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signUpInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(r.Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserNotFound {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}
