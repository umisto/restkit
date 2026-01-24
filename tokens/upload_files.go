package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UploadFilesClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope"`
}

type GenerateUploadFilesJwtRequest struct {
	SessionID uuid.UUID     `json:"session_id"`
	Issuer    string        `json:"iss"`
	Audience  []string      `json:"aud"`
	Scope     string        `json:"scope"`
	Ttl       time.Duration `json:"ttl"`
}

func NewUploadFileToken(
	req GenerateUploadFilesJwtRequest,
	sk string,
) (string, error) {
	claims := &UploadFilesClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    req.Issuer,
			Subject:   req.SessionID.String(),
			Audience:  jwt.ClaimStrings([]string{req.Issuer}),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(req.Ttl)),
		},
		Scope: req.Scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type UploadFilesJwtData struct {
	SessionID uuid.UUID
	Scope     string
}

func ParseUploadAvatarClaims(tokenStr string, sk string) (UploadFilesJwtData, error) {
	claims := UploadFilesClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})
	if err != nil || !token.Valid {
		return UploadFilesJwtData{}, err
	}

	sessionID, err := uuid.Parse(tokenStr)
	if err != nil {
		return UploadFilesJwtData{}, err
	}

	return UploadFilesJwtData{
		SessionID: sessionID,
		Scope:     claims.Scope,
	}, nil
}
