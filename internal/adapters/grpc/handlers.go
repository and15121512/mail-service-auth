package grpc

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/pkg/authgrpc"
)

func (s *Server) Validate(ctx context.Context, tokenpair *authgrpc.TokenPair) (*authgrpc.AuthResponse, error) {
	logger := s.annotatedLogger(ctx)

	tokens := &models.TokenPair{
		AuthToken:    tokenpair.AccessToken,
		RefreshToken: tokenpair.RefreshToken,
	}
	new_tokens, login, err := s.auth.ValidateAndRefresh(ctx, tokens)
	if err != nil {
		logger.Errorf("failed to validate token")
		return &authgrpc.AuthResponse{Status: "refused"}, fmt.Errorf("failed to validate token")
	}

	if *tokens != *new_tokens {
		return &authgrpc.AuthResponse{
			Status:          "refreshed",
			NewAccessToken:  new_tokens.AuthToken,
			NewRefreshToken: new_tokens.RefreshToken,
			Login:           login,
		}, nil
	}
	return &authgrpc.AuthResponse{
		Status:          "ok",
		NewAccessToken:  "",
		NewRefreshToken: "",
		Login:           login,
	}, nil
}
