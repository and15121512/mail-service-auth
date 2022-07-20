package http

import (
	"context"
	"net/http"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
)

type ctxKeyUser struct{}

func (s *Server) ValidateAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := r.Cookie("access")
			if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": "cannot read token from header",
				})
				return
			}

			user, err := s.auth.Validate(r.Context(), accessToken.Value)
			if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": "invalid token",
				})
				return
			}
			ctx := r.Context()

			ctx = context.WithValue(ctx, ctxKeyUser{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
