package http

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type Handler struct {
	authService ports.AuthService
}

func NewHandler(authService ports.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.register)
	mux.HandleFunc("POST /login", h.login)
	return mux
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		domain.User
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	createdUser, err := h.authService.Register(r.Context(), params.Username, params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: domain.User{
			ID:           createdUser.ID,
			Username:     createdUser.Username,
			PasswordHash: createdUser.PasswordHash,
			CreatedAt:    createdUser.CreatedAt,
			UpdatedAt:    createdUser.UpdatedAt,
		},
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, err := h.authService.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
