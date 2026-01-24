package mdlv

import (
	"context"
	"net/http"
	"strings"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/logium"
	"github.com/netbill/restkit/tokens"
)

const (
	AuthorizationHeader = "Authorization"
)

func AccountAuth(
	log logium.Logger,
	ctxKey int,
	sk string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(AuthorizationHeader)
			if authHeader == "" {
				log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			tokenString := parts[1]

			accountData, err := tokens.ParseAccountJWT(tokenString, sk)
			if err != nil {
				log.WithError(err).Errorf("token validation failed")
				ape.RenderErr(w, problems.Unauthorized("Token validation failed"))

				return
			}

			ctx = context.WithValue(ctx, ctxKey, tokens.AccountJwtData{
				AccountID: accountData.AccountID,
				SessionID: accountData.SessionID,
				Role:      accountData.Role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
