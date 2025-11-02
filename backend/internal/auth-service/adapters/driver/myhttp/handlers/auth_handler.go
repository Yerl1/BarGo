package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/auth-service/core/domain/dto"
	ports "backend/internal/auth-service/core/ports/driver"
	"backend/internal/mylogger"
)

type AuthHandler struct {
	mylog   mylogger.Logger
	ctx     context.Context
	service ports.IAuthService
}

func NewAuthHandler(ctx context.Context, service ports.IAuthService, mylog mylogger.Logger) *AuthHandler {
	return &AuthHandler{
		ctx:     ctx,
		service: service,
		mylog:   mylog,
	}
}

func (h *AuthHandler) SignupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle user signup
		log := h.mylog.With().Str("handler", "SignupHandler").Logger()

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user dto.SignupUser

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			JsonResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := h.service.Signup(ctx, &user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to signup user")
			JsonResponse(w, "Failed to signup user", http.StatusInternalServerError)
			return
		}
		log.Info().Msg("User signed up successfully")

		JsonResponse(w, resp, http.StatusCreated)
	}
}

func (h *AuthHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle user login
		log := h.mylog.With().Str("handler", "LoginHandler").Logger()

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user dto.LoginUser

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			JsonResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := h.service.Login(ctx, &user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to login user")
			JsonResponse(w, "Failed to login user", http.StatusInternalServerError)
			return
		}
		log.Info().Msg("User signed up successfully")

		JsonResponse(w, resp, http.StatusCreated)
	}
}

func (h *AuthHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func (h *AuthHandler) AuthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle user auth
	}
}
