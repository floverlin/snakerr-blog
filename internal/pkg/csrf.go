package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"
)

const TOKEN_NAME = "csrf_token"

func SetAndGetCSRFCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(TOKEN_NAME)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}
	token := createCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:     TOKEN_NAME,
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return token
}

func DeleteCSRFCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TOKEN_NAME,
		Value:    token,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func CheckCSRF(v url.Values, r *http.Request) (string, bool) {
	csrfCookie, err := r.Cookie(TOKEN_NAME)
	if err != nil {
		return "", false
	}
	csrfVal := v.Get(TOKEN_NAME)
	if csrfCookie.Value != csrfVal {
		return "", false
	}
	return csrfCookie.Value, true
}

func CheckAndDeleteCSRF(w http.ResponseWriter, r *http.Request, v url.Values) bool {
	if csrfCookie, ok := CheckCSRF(v, r); !ok {
		return false
	} else {
		DeleteCSRFCookie(w, csrfCookie)
		return true
	}
}

func createCSRFToken() string {
	b := make([]byte, 12)
	rand.Read(b)
	return hex.EncodeToString(b)
}
