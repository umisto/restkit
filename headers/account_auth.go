package headers

import (
	"net/http"
	"strings"

	"github.com/netbill/restkit/problems"
)

const (
	AuthorizationHeader = "Authorization"
)

func GetAuthorizationToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		return "", problems.Unauthorized("Missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", problems.Unauthorized("Missing Authorization header")
	}

	return parts[1], nil
}
