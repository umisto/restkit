package tokens

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccountClaims struct {
	jwt.RegisteredClaims
	Role      string    `json:"role"`
	SessionID uuid.UUID `json:"session_id"`
}

func (c AccountClaims) GetAccountRole() string {
	return c.Role
}

func (c AccountClaims) GetSessionID() uuid.UUID {
	return c.SessionID
}

func (c AccountClaims) GetAccountID() uuid.UUID {
	return uuid.MustParse(c.RegisteredClaims.Subject)
}

func (c AccountClaims) Validate() error {
	_, err := uuid.Parse(c.RegisteredClaims.Subject)
	if err != nil {
		return fmt.Errorf("invalid subject UUID: %w", err)
	}
	if c.SessionID == uuid.Nil {
		return fmt.Errorf("session_id cannot be nil UUID")
	}
	err = ValidateUserSystemRole(c.Role)
	if err != nil {
		return fmt.Errorf("invalid account role: %w", err)
	}

	return nil
}

func (c AccountClaims) GenerateJWT(sk string) (string, error) {
	err := c.Validate()
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(sk))
}

func ParseAccountJWT(tokenStr string, sk string) (claims AccountClaims, err error) {
	_, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(sk), nil
	})

	return claims, err
}
