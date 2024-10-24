package http

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/response"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type handler struct {
	responder         response.Responder
	userAuthenticator ports.UserAuthenticator
}

func NewHandler(responder response.Responder, userAuthenticator ports.UserAuthenticator) *handler {
	return &handler{
		userAuthenticator: userAuthenticator,
		responder:         responder,
	}
}

func (h *handler) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.register)
	mux.HandleFunc("POST /login", h.login)
	return mux
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		domain.User
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	createdUser, err := h.userAuthenticator.Register(r.Context(), params.Username, params.Password)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusCreated, response{
		User: domain.User{
			ID:           createdUser.ID,
			Username:     createdUser.Username,
			PasswordHash: createdUser.PasswordHash,
			CreatedAt:    createdUser.CreatedAt,
			UpdatedAt:    createdUser.UpdatedAt,
		},
	})
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, err := h.userAuthenticator.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
