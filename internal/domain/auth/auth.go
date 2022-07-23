package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
)

const (
	authTokenTTL             = 1 * time.Minute
	refreshTokenTTL          = 1 * time.Hour
	loginExtractionFailed    = "extracting login from token failed"
	getUserInfoFailed        = "get user info for login failed"
	invalidSignMethod        = "invalid signing method"
	tokenParsingFailed       = "token parsing failed"
	tokenClaimsParsingFailed = "token claims parsing failed"
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

type Service struct {
	db     ports.UserStorage
	logger *zap.SugaredLogger
}

func New(db ports.UserStorage, logger *zap.SugaredLogger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

// func (s *Service) Info(ctx context.Context, login string) (*models.User, error) {
// 	user, err := s.db.Get(ctx, login)
// 	if err != nil {
// 		return nil, fmt.Errorf("get user info for login %s failed: %w", login, err)
// 	}
// 	return user, nil
// }

func (s *Service) Validate(ctx context.Context, accessToken string) (*models.User, error) {
	logger := s.annotatedLogger(ctx)

	var user *models.User
	login, err := s.parseToken(ctx, accessToken)
	if err != nil {
		sentry.CaptureException(err)
		logger.Errorf(loginExtractionFailed)
		return user, fmt.Errorf(loginExtractionFailed)
	}
	user, err = s.db.Get(ctx, login)
	if err != nil {
		logger.Errorf(getUserInfoFailed)
		return user, fmt.Errorf(getUserInfoFailed)
	}
	return user, nil
}

func (s *Service) parseToken(ctx context.Context, accessToken string) (string, error) {
	logger := s.annotatedLogger(ctx)

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Errorf(invalidSignMethod)
			return nil, fmt.Errorf(invalidSignMethod)
		}
		return []byte(config.GetConfig(logger).Auth.Secret), nil
	})
	if err != nil {
		logger.Errorf(tokenParsingFailed)
		return "", fmt.Errorf(tokenParsingFailed)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		logger.Errorf(tokenClaimsParsingFailed)
		return "", fmt.Errorf(tokenClaimsParsingFailed)
	}
	return claims.Login, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (models.TokenPair, error) {
	logger := s.annotatedLogger(ctx)

	var tokens models.TokenPair
	userModel, err := s.db.Get(ctx, login)
	if err != nil {
		logger.Errorf("get user info for login %s failed", login)
		return tokens, fmt.Errorf("get user info for login %s failed", login)
	}

	passwordHash := s.generatePasswordHash(ctx, password)
	if userModel.PasswordHash != passwordHash {
		logger.Errorf("invalid password for login %s", login)
		return tokens, fmt.Errorf("invalid password for login %s", login)
	}

	authToken, err := s.generateToken(ctx, login, authTokenTTL)
	if err != nil {
		logger.Errorf("generate auth token for login %s failed", login)
		return tokens, fmt.Errorf("generate auth token for login %s failed", login)
	}
	refreshToken, err := s.generateToken(ctx, login, refreshTokenTTL)
	if err != nil {
		logger.Errorf("generate refresh token for login %s failed", login)
		return tokens, fmt.Errorf("generate refresh token for login %s failed", login)
	}
	tokens = models.TokenPair{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func (s *Service) generatePasswordHash(ctx context.Context, password string) string {
	logger := s.annotatedLogger(ctx)

	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger).Auth.Salt)))
}

func (s *Service) generateToken(ctx context.Context, login string, tokenTTL time.Duration) (string, error) {
	logger := s.annotatedLogger(ctx)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		login,
	})
	return token.SignedString([]byte(config.GetConfig(logger).Auth.Secret))
}
