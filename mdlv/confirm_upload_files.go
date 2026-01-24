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

const UploadHeader = "Upload-Token"

func ConfirmUploadFiles(
	log logium.Logger,
	ctxKey int,
	sk string,
	scope string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(UploadHeader)
			if authHeader == "" {
				log.Errorf("missing %s header", UploadHeader)
				ape.RenderErr(w, problems.Unauthorized("Missing Upload-Token header"))

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			tokenString := parts[1]

			uploadSessionData, err := tokens.ParseUploadAvatarClaims(tokenString, sk)
			if err != nil {
				log.WithError(err).Errorf("token validation failed")
				ape.RenderErr(w, problems.Unauthorized("Token validation failed"))

				return
			}

			if uploadSessionData.Scope != scope {
				log.Errorf("invalid token scope")
				ape.RenderErr(w, problems.Unauthorized("Invalid token scope"))

				return
			}

			ctx = context.WithValue(ctx, ctxKey, uploadSessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
