package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
)

func (s *Server) authHandlers() http.Handler {
	h := chi.NewRouter()
	h.With(s.ValidateAuth()).Get("/i", s.Info)
	h.Post("/login", s.Login)
	h.Post("/logout", s.Logout)
	return h
}

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(ctxKeyUser{}).(*models.User)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to extract token from request",
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]string{
		"login": user.Login,
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid Authorization header",
		})
		return
	}
	tokens, err := s.auth.Login(r.Context(), user, password)
	if err != nil {
		utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
			"error": err.Error(),
		})
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
