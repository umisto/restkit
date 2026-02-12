package grants

import (
	"net/http"
	"strings"

	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/tokens"
)

const UploadHeader = "Upload-Token"

type UploadContentParams struct {
	Audience   string
	Resource   string
	ResourceID string
}

func UploadContentGrant(
	// request is the HTTP request containing the Upload-Token header
	r *http.Request,
	// uploadSK is the secret key used to validate the token
	uploadSK string,
	// params contains the expected parameters of the token
	params UploadContentParams,
) (tokens.UploadContentClaims, error) {
	authHeader := r.Header.Get(UploadHeader)
	if authHeader == "" {
		return tokens.UploadContentClaims{}, problems.Unauthorized("Missing Upload-Token header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return tokens.UploadContentClaims{}, problems.Unauthorized("Missing Upload-Token header")
	}

	tokenString := parts[1]

	uploadSessionData, err := tokens.ParseUploadFilesClaims(tokenString, uploadSK)
	if err != nil {
		return tokens.UploadContentClaims{}, problems.Unauthorized("upload token validation failed")
	}

	if uploadSessionData.ResourceType != params.Resource {
		return tokens.UploadContentClaims{}, problems.Unauthorized("invalid upload token resource")
	}

	if uploadSessionData.ResourceID != params.ResourceID {
		return tokens.UploadContentClaims{}, problems.Unauthorized("invalid upload token resource ID")
	}

	audSuccess := false
	for _, aud := range uploadSessionData.Audience {
		if aud == params.Audience {
			audSuccess = true
			break
		}
	}
	if !audSuccess {
		return tokens.UploadContentClaims{}, problems.Unauthorized("invalid upload token audience")
	}

	return uploadSessionData, nil
}
