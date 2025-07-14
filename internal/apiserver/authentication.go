package apiserver

import (
	"blog/internal/model"
	"blog/internal/pkg"
	"blog/internal/storage"
	"blog/internal/templates"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int8
type role int8

const (
	ctxIDKey   ctxKey = iota // user id uint64
	ctxRoleKey               // user role type=role
	ctxUserKey               // *model.User finded by id
)

const (
	userRole role = iota
)

func (s *APIServer) handleRegister() http.HandlerFunc {
	const op = "handle register"
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			vals := r.Form
			if !(vals.Has("email") && vals.Has("username") && vals.Has("password")) {
				s.errorPageReferer(w, r, errors.New("some values are empty"))
				return
			}
			if !pkg.CheckAndDeleteCSRF(w, r, vals) {
				s.errorPageReferer(w, r, errors.New("wrong csrf token"))
				return
			}
			u := model.User{
				Email:    vals.Get("email"),
				Username: vals.Get("username"),
				Password: vals.Get("password"),
			}
			if err := u.ValidateBeforeCreate(); err != nil {
				s.errorPageReferer(w, r, err)
				return
			}
			u.Password = pkg.HashPassword(u.Password, s.config.Secret)
			if _, err := s.storage.User().Create(r.Context(), &u); err != nil {
				if errors.Is(err, storage.ErrUnique) {
					s.errorPageReferer(w, r, errors.New("email is already zan'yat"))
					return
				} else {
					s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
					return
				}
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			csrfToken := pkg.SetAndGetCSRFCookie(w, r)
			exeData := templates.RegisterPageData(r.Context(), "Register", csrfToken)
			s.renderPage(w, r, http.StatusOK, "register", exeData)
		}
	}
}

func (s *APIServer) handleLogin() http.HandlerFunc {
	const op = "handle login"
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			vals := r.Form
			if !(vals.Has("email") && vals.Has("password")) {
				s.errorPageReferer(w, r, errors.New("some values are empty"))
				return
			}
			if !pkg.CheckAndDeleteCSRF(w, r, vals) {
				s.errorPageReferer(w, r, errors.New("wrong csrf token"))
				return
			}
			u, err := s.storage.User().FindByEmail(r.Context(), vals.Get("email"))
			if err != nil {
				if errors.Is(err, storage.ErrNoRows) {
					s.errorPageReferer(w, r, errors.New("email not found"))
					return
				} else {
					s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
					return
				}
			}
			if u.Password != pkg.HashPassword(vals.Get("password"), s.config.Secret) {
				s.errorPageReferer(w, r, errors.New("wrong password"))
				return
			}
			tokenString, err := generateUserToken(u.ID, userRole, s.config.Secret)
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			cookie := &http.Cookie{
				Name:     "token",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Value:    tokenString,
				Expires:  time.Now().Add(time.Hour * 24 * 60),
				MaxAge:   int(time.Hour * 24 * 60 / time.Second),
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/my_page", http.StatusSeeOther)
		} else {
			csrfToken := pkg.SetAndGetCSRFCookie(w, r)
			exeData := templates.LoginPageData(r.Context(), "Login", csrfToken)
			s.renderPage(w, r, http.StatusOK, "login", exeData)
		}
	}
}

func (s *APIServer) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			HttpOnly: true,
			Name:     "token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *APIServer) authMWare(next http.Handler) http.Handler {
	const op = "auth middleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			next.ServeHTTP(w, r)
			return
		}
		if err := cookie.Valid(); err != nil {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := cookie.Value
		token, err := parseToken(tokenString, s.config.Secret)
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := uint64(claims["id"].(float64))
		userRole := role(claims["role"].(float64))
		u, err := s.storage.User().FindByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRows) {
				next.ServeHTTP(w, r)
				return
			} else {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxIDKey, userID)
		ctx = context.WithValue(ctx, ctxUserKey, u)
		ctx = context.WithValue(ctx, ctxRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *APIServer) authSecureMWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLogged(r) {
			s.setFlash(w, "please, login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseToken(tokenString string, key string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
}

func generateUserToken(id uint64, role role, key string) (string, error) {
	liveTime := time.Now().Add(48 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  liveTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func getCurrentID(r *http.Request) uint64 {
	v := r.Context().Value(ctxIDKey)
	if v == nil {
		return 0
	}
	return v.(uint64)
}

func isLogged(r *http.Request) bool {
	return r.Context().Value(ctxIDKey) != nil
}
