package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccountClaims struct {
	jwt.RegisteredClaims
	Role      string    `json:"role"`
	SessionID uuid.UUID `json:"session_id"`
}

type GenerateAccountJwtRequest struct {
	Issuer    string        `json:"iss"`
	AccountID uuid.UUID     `json:"sub"`
	SessionID uuid.UUID     `json:"session_id"`
	Role      string        `json:"role"`
	Ttl       time.Duration `json:"ttl"`
}

func GenerateAccountJWT(
	req GenerateAccountJwtRequest,
	sk string,
) (string, error) {
	claims := &AccountClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    req.Issuer,
			Subject:   req.AccountID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(req.Ttl)),
		},
		SessionID: req.SessionID,
		Role:      req.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type AccountJwtData struct {
	AccountID uuid.UUID
	SessionID uuid.UUID
	Role      string
}

func ParseAccountJWT(tokenStr string, sk string) (AccountJwtData, error) {
	claims := AccountClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(sk), nil
	})

	if err != nil || !token.Valid {
		return AccountJwtData{}, err
	}

	accountID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return AccountJwtData{}, err
	}

	return AccountJwtData{
		AccountID: accountID,
		SessionID: claims.SessionID,
		Role:      claims.Role,
	}, nil
}
