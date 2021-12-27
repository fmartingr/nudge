package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fmartingr/nudge/internal/http"
	"github.com/fmartingr/nudge/internal/pinger"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http   *http.HttpServer
	pinger *pinger.Pinger
	logger *logrus.Entry
}

func (s *Server) Start(ctx context.Context) {
	go s.http.Start()
	go s.pinger.Start(ctx)
}

func (s *Server) Stop() error {
	if err := s.http.Stop(); err != nil {
		s.logger.Errorf("error stopping http server: %s", err)
	}

	s.pinger.Stop()

	return nil
}

func (s *Server) WaitStop() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	return s.Stop()
}

func NewServer(conf *Config) *Server {
	loglevel, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		logrus.Warnf("Error parsing loglevel '%s', using default.", loglevel)
	} else {
		logrus.SetLevel(loglevel)
	}

	logger := logrus.WithField("service", "nudge")
	ping := pinger.NewPinger(logger, conf.IPs, conf.Interval)

	return &Server{
		logger: logger.WithField("from", "server"),
		pinger: ping,
		http:   http.NewHttpServer(logger, conf.Port, ping),
	}
}
