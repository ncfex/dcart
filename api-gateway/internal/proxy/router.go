package proxy

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ncfex/dcart/api-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart/api-gateway/internal/middleware"
)

type Router struct {
	cfg      *config.Config
	services map[string]*serviceProxy
	auth     *middleware.AuthMiddleware
	mu       sync.RWMutex
}

func NewRouter(cfg *config.Config) (*Router, error) {
	r := &Router{
		cfg:      cfg,
		services: make(map[string]*serviceProxy),
		auth:     middleware.NewAuthMiddleware(cfg.Auth),
	}

	if err := r.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}
	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(req.URL.Path, "/"), "/", 2)
	if len(parts) == 0 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	serviceName := parts[0]

	r.mu.RLock()
	service, exists := r.services[serviceName]
	r.mu.RUnlock()

	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if len(parts) > 1 {
		req.URL.Path = "/" + parts[1]
	} else {
		req.URL.Path = "/"
	}

	service.ServeHTTP(w, req)
}

func (r *Router) initializeServices() error {
	for _, svc := range r.cfg.Services {
		proxy, err := newServiceProxy(&svc, r)
		if err != nil {
			return fmt.Errorf("failed to initialize service %s: %w", svc.Name, err)
		}

		r.mu.Lock()
		r.services[svc.Name] = proxy
		r.mu.Unlock()
	}
	return nil
}
