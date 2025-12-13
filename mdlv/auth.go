package mdlv

import (
	"context"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/restkit/token"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "Authorization"
)

func Auth(ctxKey interface{}, skUser string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(AuthorizationHeader)
			if authHeader == "" {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			tokenString := parts[1]

			userData, err := token.VerifyAccountJWT(tokenString, skUser)
			if err != nil {
				ape.RenderErr(w,
					problems.Unauthorized("Token validation failed"),
				)

				return
			}

			userID, err := uuid.Parse(userData.Subject)
			if err != nil {
				ape.RenderErr(w,
					problems.Unauthorized("UserID ID is nov valid"),
				)

				return
			}

			ctx = context.WithValue(ctx, ctxKey, token.AccountData{
				ID:        userID,
				SessionID: userData.SessionID,
				Role:      userData.Role,
				Username:  userData.Username,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
