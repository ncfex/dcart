package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ncfex/dcart/api-gateway/internal/infrastructure/config"
)

type serviceProxy struct {
	cfg     *config.ServiceConfig
	proxy   *httputil.ReverseProxy
	handler http.Handler
}

func newServiceProxy(cfg *config.ServiceConfig, router *Router) (*serviceProxy, error) {
	target, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	var handler http.Handler = proxy

	if cfg.RequiresAuth {
		handler = router.auth.Middleware(handler)
	}

	return &serviceProxy{
		cfg:     cfg,
		proxy:   proxy,
		handler: handler,
	}, nil
}

func (sp *serviceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), sp.cfg.Timeout)
	defer cancel()

	sp.handler.ServeHTTP(w, r.WithContext(ctx))
}
