package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/middleware"
	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/request"
	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/response"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type handler struct {
	responder         response.Responder
	userAuthenticator ports.UserAuthenticator
	tokenManager      ports.TokenManager
	tokenRepo         ports.TokenRepository
	userRepo          ports.UserRepository
}

func NewHandler(
	responder response.Responder,
	userAuthenticator ports.UserAuthenticator,
	tokenManager ports.TokenManager,
	tokenRepo ports.TokenRepository,
	userRepo ports.UserRepository,
) *handler {
	return &handler{
		userAuthenticator: userAuthenticator,
		responder:         responder,
		tokenManager:      tokenManager,
		tokenRepo:         tokenRepo,
		userRepo:          userRepo,
	}
}

func (h *handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	publicChain := middleware.NewChain(
		middleware.Recovery(h.responder),
		middleware.Logger(),
	)

	refreshTokenRequiredChain := middleware.NewChain(
		middleware.Recovery(h.responder),
		middleware.Logger(),
		middleware.AuthenticateWithRefreshToken(h.tokenManager, h.tokenRepo, h.userRepo),
	)

	accessTokenProtectedChain := middleware.NewChain(
		middleware.Recovery(h.responder),
		middleware.Logger(),
		middleware.AuthenticateWithJWT(h.tokenManager, h.tokenRepo, h.userRepo),
	)

	// public
	mux.Handle("POST /register", publicChain.ThenFunc(h.register))
	mux.Handle("POST /login", publicChain.ThenFunc(h.login))

	// protected
	mux.Handle("GET /profile", accessTokenProtectedChain.ThenFunc(h.profile))

	// refresh required
	mux.Handle("POST /refresh", refreshTokenRequiredChain.ThenFunc(h.refresh))
	mux.Handle("POST /logout", refreshTokenRequiredChain.ThenFunc(h.logout))

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
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	tokenPair, err := h.userAuthenticator.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		Token:        string(tokenPair.AccessToken),
		RefreshToken: string(tokenPair.RefreshToken),
	})
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	err = h.userAuthenticator.Logout(r.Context(), refreshToken)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	tokenPair, err := h.userAuthenticator.Refresh(r.Context(), refreshToken)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		Token: string(tokenPair.AccessToken),
	})
}

func (h *handler) profile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User domain.User `json:"user"`
	}

	user, exists := request.GetUserFromContext(r.Context())
	if !exists {
		h.responder.RespondWithError(w, http.StatusNotFound, "no user found", errors.New("no user found"))
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		User: domain.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
