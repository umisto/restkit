package tokens

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UploadContentClaims struct {
	jwt.RegisteredClaims
	UploadSessionID uuid.UUID `json:"upload_session_id"`
	ResourceID      string    `json:"resource_id"`
	ResourceType    string    `json:"resourcetype"`
}

func (c UploadContentClaims) GetOwnerAccountID() uuid.UUID {
	return uuid.MustParse(c.RegisteredClaims.Subject)
}

func (c UploadContentClaims) GetSessionID() uuid.UUID {
	return c.UploadSessionID
}

func (c UploadContentClaims) GetResourceID() string {
	return c.ResourceID
}

func (c UploadContentClaims) GetResource() string {
	return c.ResourceType
}

func (c UploadContentClaims) Validate() error {
	_, err := uuid.Parse(c.RegisteredClaims.Subject)
	if err != nil {
		return fmt.Errorf("invalid subject UUID: %w", err)
	}
	if c.UploadSessionID == uuid.Nil {
		return fmt.Errorf("upload_session_id cannot be nil UUID")
	}
	if c.ResourceID == "" {
		return fmt.Errorf("resource_id cannot be empty")
	}
	if c.ResourceType == "" {
		return fmt.Errorf("resource_type cannot be empty")
	}

	return nil
}

func (c UploadContentClaims) GenerateJWT(sk string) (string, error) {
	err := c.Validate()
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(sk))
}

func ParseUploadFilesClaims(tokenStr string, sk string) (claims UploadContentClaims, err error) {
	_, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(sk), nil
	})
	return claims, err
}

type UploadContent interface {
	GetOwnerAccountID() uuid.UUID
	GetSessionID() uuid.UUID
	GetResourceID() string
	GetResource() string
}
