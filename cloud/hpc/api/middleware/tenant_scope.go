package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TenantScopeMiddleware struct{}

func NewTenantScopeMiddleware() *TenantScopeMiddleware {
	return &TenantScopeMiddleware{}
}

func (m *TenantScopeMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Context().Value("auth")
		if auth == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tenantFromToken := auth.(*AuthContext).TenantID
		tenantFromPath := chi.URLParam(r, "tenant")

		if tenantFromPath != "" && tenantFromPath != tenantFromToken {
			http.Error(w, "tenant access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
