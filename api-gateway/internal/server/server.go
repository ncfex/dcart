package server

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart/api-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart/api-gateway/internal/proxy"
)

type Server struct {
	cfg    *config.Config
	log    *log.Logger
	router *proxy.Router
	Server *http.Server
}

func NewServer(cfg *config.Config, log *log.Logger) (*Server, error) {
	router, err := proxy.NewRouter(cfg)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	return &Server{
		cfg:    cfg,
		log:    log,
		router: router,
		Server: srv,
	}, nil
}
