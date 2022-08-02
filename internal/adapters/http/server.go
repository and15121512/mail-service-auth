package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/ports"
	"go.uber.org/zap"
)

type Server struct {
	auth   ports.Auth
	server *http.Server
	l      net.Listener
	port   int
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger, auth ports.Auth) (*Server, error) {
	var (
		err error
		s   Server
	)
	s.l, err = net.Listen("tcp", ":"+config.GetConfig(logger).Ports.HttpPort)
	if err != nil {
		logger.Fatalf("Failed listen port: %s", err)
	}
	s.auth = auth
	s.port = s.l.Addr().(*net.TCPAddr).Port
	s.server = &http.Server{
		Handler: s.routes(),
		// Addr   : ":" + strconv.Itoa(s.port),  # server.Serve(s.l) does this shit!
	}
	s.logger = logger

	return &s, nil
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.l); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/healthz", s.healthzHandler)
	r.Mount("/", s.authHandlers())
	r.Mount("/debug/", middleware.Profiler())

	return r
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
