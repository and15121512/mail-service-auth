package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
	"go.uber.org/zap"
)

const (
	tokenExtractionFailed = "failed to extract token from request"
	invalidAuthHeader     = "invalid Authorization header"
)

func (s *Server) authHandlers() http.Handler {
	h := chi.NewRouter()
	h.With(s.AnnotateContext()).With(s.ValidateAuth()).Get("/i", s.Info)
	h.With(s.AnnotateContext()).Post("/login", s.Login)
	h.With(s.AnnotateContext()).Post("/logout", s.Logout)
	return h
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

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	user, ok := r.Context().Value(ctxKeyUser{}).(*models.User)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": tokenExtractionFailed,
		})
		logger.Errorf(tokenExtractionFailed)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]string{
		"login": user.Login,
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	user, password, ok := r.BasicAuth()
	if !ok {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": invalidAuthHeader,
		})
		logger.Errorf(invalidAuthHeader)
		return
	}
	tokens, err := s.auth.Login(r.Context(), user, password)
	if err != nil {
		utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
			"error": err.Error(),
		})
		logger.Errorf(err.Error())
		return
	}
	redirectURI := r.URL.Query().Get("redirect_uri")
	status := http.StatusOK
	if redirectURI != "" {
		http.Redirect(w, r, redirectURI, http.StatusSeeOther)
		status = http.StatusSeeOther
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access",
		Value: tokens.AuthToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh",
		Value: tokens.RefreshToken,
	})
	utils.ResponseJSON(w, status, map[string]string{
		"accessToken":  tokens.AuthToken,
		"refreshToken": tokens.RefreshToken,
	})
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	utils.ResponseJSON(w, http.StatusOK, map[string]string{
		"accessToken":  "",
		"refreshToken": "",
	})
}
