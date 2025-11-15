package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UsersClaims struct {
	jwt.RegisteredClaims
	Role      string    `json:"role"`
	SessionID uuid.UUID `json:"session_id"`
}

func VerifyUserJWT(tokenString, sk string) (UsersClaims, error) {
	claims := UsersClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})
	if err != nil || !token.Valid {
		return UsersClaims{}, err
	}
	return claims, nil
}

type GenerateUserJwtRequest struct {
	Issuer    string        `json:"iss"`
	Audience  []string      `json:"aud"`
	UserID    uuid.UUID     `json:"sub"`
	SessionID uuid.UUID     `json:"session_id"`
	Role      string        `json:"role"`
	Ttl       time.Duration `json:"ttl"`
}

func GenerateUserJWT(
	request GenerateUserJwtRequest,
	sk string,
) (string, error) {
	expirationTime := time.Now().UTC().Add(request.Ttl * time.Second)
	claims := &UsersClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    request.Issuer,
			Subject:   request.UserID.String(),
			Audience:  jwt.ClaimStrings(request.Audience),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		SessionID: request.SessionID,
		Role:      request.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type UserData struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	Role      string
}
