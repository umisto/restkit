package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccountClaims struct {
	jwt.RegisteredClaims
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	SessionID uuid.UUID `json:"session_id"`
}

func VerifyAccountJWT(tokenString, sk string) (AccountClaims, error) {
	claims := AccountClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})
	if err != nil || !token.Valid {
		return AccountClaims{}, err
	}
	return claims, nil
}

type GenerateAccountJwtRequest struct {
	Issuer    string        `json:"iss"`
	Audience  []string      `json:"aud"`
	AccountID uuid.UUID     `json:"sub"`
	SessionID uuid.UUID     `json:"session_id"`
	Role      string        `json:"role"`
	Username  string        `json:"username"`
	Ttl       time.Duration `json:"ttl"`
}

func GenerateAccountJWT(
	request GenerateAccountJwtRequest,
	sk string,
) (string, error) {
	expirationTime := time.Now().UTC().Add(request.Ttl * time.Second)
	claims := &AccountClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    request.Issuer,
			Subject:   request.AccountID.String(),
			Audience:  jwt.ClaimStrings(request.Audience),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Username:  request.Username,
		SessionID: request.SessionID,
		Role:      request.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type AccountData struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	Role      string
	Username  string
}
