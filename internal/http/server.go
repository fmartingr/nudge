package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/fmartingr/nudge/internal/pinger"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	server *http.Server
	logger *logrus.Entry
}

func (s *HttpServer) Start() {
	s.logger.Infof("Starting HTTP server on %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Errorf("Error starting HTTP Server: %s", err)
	}
}

func (s *HttpServer) Stop() error {
	shuwdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(shuwdownContext); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Errorf("Error stopping HTTP Server: %s", err)
	}

	return nil
}

func NewHttpServer(logger *logrus.Entry, port int, ping *pinger.Pinger) *HttpServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/status", statusHandler(ping))

	return &HttpServer{
		logger: logger.WithField("from", "http"),
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
	}
}
