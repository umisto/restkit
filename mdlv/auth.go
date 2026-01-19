package mdlv

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
	"github.com/netbill/restkit/auth"
)

const (
	AuthorizationHeader = "Authorization"
)

func (s Service) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(AuthorizationHeader)
			if authHeader == "" {
				s.log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				s.log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			tokenString := parts[1]

			userData, err := auth.VerifyAccountJWT(tokenString, s.skUser)
			if err != nil {
				s.log.WithError(err).Errorf("token validation failed")
				ape.RenderErr(w, problems.Unauthorized("Token validation failed"))

				return
			}

			userID, err := uuid.Parse(userData.Subject)
			if err != nil {
				s.log.WithError(err).Errorf("subject filed in token is nov valid")
				ape.RenderErr(w, problems.Unauthorized("Token validation failed"))

				return
			}

			ctx = context.WithValue(ctx, s.ctxKey, auth.AccountData{
				ID:        userID,
				SessionID: userData.SessionID,
				Role:      userData.Role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
