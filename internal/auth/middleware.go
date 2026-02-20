package auth

import (
	"context"
	"net/http"
	"strings"

	"zenx/internal/router"
)

type authCtxKey string

const claimsKey authCtxKey = "auth_claims"

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}

func JWTAuth(mgr *JWTManager) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			authz := c.Request.Header.Get("Authorization")
			if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
				http.Error(c.Writer, "missing bearer token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authz, "Bearer ")
			claims, err := mgr.Parse(token)
			if err != nil {
				http.Error(c.Writer, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(c.Request.Context(), claimsKey, claims)
			c.WithContext(ctx)
			next(c)
		}
	}
}

func RequireRoles(roles ...string) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			claims, ok := ClaimsFromContext(c.Request.Context())
			if !ok || !HasRole(claims.Roles, roles...) {
				http.Error(c.Writer, "forbidden", http.StatusForbidden)
				return
			}
			next(c)
		}
	}
}
