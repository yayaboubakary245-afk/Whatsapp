package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/groupe12-supptic/whatsapp-campus/internal/auth"
	"github.com/groupe12-supptic/whatsapp-campus/internal/security"
	"github.com/groupe12-supptic/whatsapp-campus/pkg/utils"
)

type ctxKey string

const (
	CtxUserID ctxKey = "user_id"
	CtxEmail  ctxKey = "email"
	CtxRole   ctxKey = "role"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(rl *security.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow(r.RemoteAddr) {
				utils.WriteError(w, http.StatusTooManyRequests, "Trop de requêtes, veuillez patienter")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenStr := ""
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			} else if q := r.URL.Query().Get("token"); q != "" {
				tokenStr = q
			}
			if tokenStr == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Jeton JWT manquant")
				return
			}
			claims, err := auth.ParseJWT(jwtSecret, tokenStr)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "Jeton JWT invalide ou expiré")
				return
			}
			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(CtxRole).(string)
		if role != "admin" {
			utils.WriteError(w, http.StatusForbidden, "Accès réservé aux administrateurs")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserIDFromCtx(r *http.Request) string {
	v, _ := r.Context().Value(CtxUserID).(string)
	return v
}

func RoleFromCtx(r *http.Request) string {
	v, _ := r.Context().Value(CtxRole).(string)
	return v
}
