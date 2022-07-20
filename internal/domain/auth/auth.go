package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/ports"
)

const (
	authTokenTTL    = 1 * time.Minute
	refreshTokenTTL = 1 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

type Service struct {
	db ports.UserStorage
}

func New(db ports.UserStorage) *Service {
	return &Service{
		db: db,
	}
}

// func (s *Service) Info(ctx context.Context, login string) (*models.User, error) {
// 	user, err := s.db.Get(ctx, login)
// 	if err != nil {
// 		return nil, fmt.Errorf("get user info for login %s failed: %w", login, err)
// 	}
// 	return user, nil
// }

func (s *Service) Validate(ctx context.Context, accessToken string) (*models.User, error) {
	var user *models.User
	login, err := parseToken(accessToken)
	if err != nil {
		return user, fmt.Errorf("extracting login from token failed: %w", err)
	}
	user, err = s.db.Get(ctx, login)
	if err != nil {
		return user, fmt.Errorf("get user info for login %s failed: %w", login, err)
	}
	return user, nil
}

func parseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(config.GetConfig().Auth.Secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("token claims parsing failed")
	}
	return claims.Login, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (models.TokenPair, error) {
	var tokens models.TokenPair
	userModel, err := s.db.Get(ctx, login)
	if err != nil {
		return tokens, fmt.Errorf("get user info for login %s failed: %w", login, err)
	}

	passwordHash := generatePasswordHash(password)
	if userModel.PasswordHash != passwordHash {
		return tokens, fmt.Errorf("invalid password for login %s", login)
	}

	authToken, err := generateToken(login, authTokenTTL)
	if err != nil {
		return tokens, fmt.Errorf("generate auth token for login %s failed", login)
	}
	refreshToken, err := generateToken(login, refreshTokenTTL)
	if err != nil {
		return tokens, fmt.Errorf("generate refresh token for login %s failed", login)
	}
	tokens = models.TokenPair{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig().Auth.Salt)))
}

func generateToken(login string, tokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		login,
	})
	return token.SignedString([]byte(config.GetConfig().Auth.Secret))
}
