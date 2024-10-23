package http

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
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
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
