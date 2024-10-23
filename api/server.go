package api

import (
	"context"
	"errors"
	"net/http"
	"time"
	"twitter/api/handler"
	"twitter/service"
	"twitter/storage"

	"github.com/gin-gonic/gin"
	"twitter/config"
	"twitter/pkg/logger"
)

type Server struct {
	log        logger.Logger
	cfg        config.Config
	services   service.IServiceManager
	storage    storage.IStorage
	router     *gin.Engine
	httpServer *http.Server
	handler    handler.Handler
}

func New(cfg config.Config, services service.IServiceManager, storage storage.IStorage, log logger.Logger, h handler.Handler) *Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	s := &Server{
		cfg:      cfg,
		services: services,
		storage:  storage,
		log:      log,
		handler:  h,
		router:   r,
	}

	s.endpoints()

	s.httpServer = &http.Server{
		Addr:              cfg.HTTPort,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 10,
	}

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
