package grants

import (
	"net/http"
	"strings"

	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/tokens"
)

const (
	AuthorizationHeader = "Authorization"
)

func AccountAuthToken(
	// request is the HTTP request containing the Authorization header
	r *http.Request,
	// accessSK is the secret key used to validate the token
	accessSK string,
	// audience is the expected audience of the token
	audience string,
	// roles is the list of allowed roles for the token
	roles ...string,
) (tokens.AccountClaims, error) {
	authHeader := r.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		return tokens.AccountClaims{}, problems.Unauthorized("Missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return tokens.AccountClaims{}, problems.Unauthorized("Missing Authorization header")
	}

	tokenString := parts[1]

	accountData, err := tokens.ParseAccountJWT(tokenString, accessSK)
	if err != nil {
		return tokens.AccountClaims{}, problems.Unauthorized("Token validation failed")
	}

	if err = tokens.ValidateUserSystemRole(accountData.Role); err != nil {
		return tokens.AccountClaims{}, problems.Unauthorized("account role not valid")
	}

	if audience != "" {
		access := false
		for _, aud := range accountData.Audience {
			if aud == audience {
				access = true
				break
			}
		}
		if !access {
			return tokens.AccountClaims{}, problems.Unauthorized("invalid token audience")
		}
	}

	if roles != nil || len(roles) > 0 {
		allowed := false
		for _, role := range roles {
			if accountData.Role == role {
				allowed = true
				break
			}
		}
		if !allowed {
			return tokens.AccountClaims{}, problems.Forbidden("account role not allowed")
		}
	}

	return accountData, nil
}
