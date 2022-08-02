package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
	"gitlab.com/sukharnikov.aa/mail-service-auth/pkg/authgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	authgrpc.UnimplementedAuthGrpcServer
	auth   ports.Auth
	server *grpc.Server
	l      net.Listener
	port   int
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger, auth ports.Auth) (*Server, error) {
	var (
		s   Server
		err error
	)

	s.l, err = net.Listen("tcp", ":"+config.GetConfig(logger).Ports.GrpcPort)
	if err != nil {
		logger.Fatalf("failed listen port: %s", err)
	}
	s.auth = auth
	s.port = s.l.Addr().(*net.TCPAddr).Port

	s.server = grpc.NewServer()
	s.logger = logger
	authgrpc.RegisterAuthGrpcServer(s.server, &s)

	return &s, nil
}

func (s *Server) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
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
	s.server.Stop()
	return nil
}
