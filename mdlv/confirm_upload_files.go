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

type ConfirmUploadFilesParams struct {
	Audience   string
	Resource   string
	ResourceID string
}

func ConfirmUploadFiles(log *logium.Logger, ctxKey int, sk string, params ConfirmUploadFilesParams) func(http.Handler) http.Handler {
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
				log.Errorf("missing Upload-Token header")
				ape.RenderErr(w, problems.Unauthorized("Missing Upload-Token header"))

				return
			}

			tokenString := parts[1]

			uploadSessionData, err := tokens.ParseUploadFilesClaims(tokenString, sk)
			if err != nil {
				log.WithError(err).Errorf("upload token validation failed")
				ape.RenderErr(w, problems.Unauthorized("upload token validation failed"))

				return
			}

			if uploadSessionData.Resource != params.Resource {
				log.Errorf("invalid upload token resource")
				ape.RenderErr(w, problems.Unauthorized("invalid upload token resource"))

				return
			}

			if uploadSessionData.ResourceID != params.ResourceID {
				log.Errorf("invalid upload token resource id")
				ape.RenderErr(w, problems.Unauthorized("invalid upload token resource id"))

				return
			}

			audSuccess := false
			for _, aud := range uploadSessionData.Audience {
				if aud == params.Audience {
					audSuccess = true
					break
				}
			}
			if !audSuccess {
				log.Errorf("invalid upload token audience")
				ape.RenderErr(w, problems.Unauthorized("invalid upload token audience"))

				return
			}

			ctx = context.WithValue(ctx, ctxKey, uploadSessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
