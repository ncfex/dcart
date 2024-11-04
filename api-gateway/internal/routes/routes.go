package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Router struct {
	authService *url.URL
}

func NewRouter(authURL string) (*Router, error) {
	auth, err := url.Parse(authURL)
	if err != nil {
		return nil, err
	}

	return &Router{
		authService: auth,
	}, nil
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// auth service proxy
	mux.Handle("/auth/", http.StripPrefix("/auth", httputil.NewSingleHostReverseProxy(r.authService)))

	return mux
}
