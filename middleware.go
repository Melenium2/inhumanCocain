package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	auth "github.com/inhumanLightBackend/auth/logic"
	gates "github.com/inhumanLightBackend/cocain/logic"
)

var (
	errNotAuthorized = errors.New("Not authorized")
	except = []string{ "/auth", "/user/create" }
)

type Middleware interface {
	Translate(http.Handler) http.Handler
}

type middleware struct {
	service gates.Gates
}

func NewMiddleware(s gates.Gates) Middleware {
	return &middleware {
		service: s,
	}
}

func (m *middleware) Translate(next http.Handler) http.Handler {	
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if exceptPattern(r.RequestURI) {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("X-Opaque")
		r.Header.Del("X-Opaque")
		if token == "" {
			SendError(w, r, http.StatusUnauthorized, errNotAuthorized)
			return
		}
		jwt, err := m.service.Translate(r.Context(), token)
		if err != nil {
			SendError(w, r, http.StatusUnauthorized, errNotAuthorized)
			return
		}
		c := context.WithValue(r.Context(), auth.CtxUserKey, jwt)
		r.Header.Set("Authorization", jwt)
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func exceptPattern(uri string) bool {
	for _, path := range except {
		if strings.Contains(uri, path) {
			return true
		}
	}
	return false
}