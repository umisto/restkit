package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UploadFilesClaims struct {
	jwt.RegisteredClaims
	UploadSessionID uuid.UUID `json:"upload_session_id"`
	ResourceID      string    `json:"resource_id"`
	Resource        string    `json:"resource"`
}

type GenerateUploadFilesJwtRequest struct {
	OwnerAccountID  uuid.UUID     `json:"sub"`
	UploadSessionID uuid.UUID     `json:"upload_session_id"`
	ResourceID      string        `json:"resource_id"`
	Resource        string        `json:"resource"`
	Issuer          string        `json:"iss"`
	Audience        []string      `json:"aud"`
	Ttl             time.Duration `json:"ttl"`
}

func NewUploadFileToken(
	req GenerateUploadFilesJwtRequest,
	sk string,
) (string, error) {
	claims := &UploadFilesClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    req.Issuer,
			Subject:   req.OwnerAccountID.String(),
			Audience:  jwt.ClaimStrings(req.Audience),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(req.Ttl)),
		},
		UploadSessionID: req.UploadSessionID,
		ResourceID:      req.ResourceID,
		Resource:        req.Resource,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type UploadFilesJwtData struct {
	OwnerAccountID  uuid.UUID
	UploadSessionID uuid.UUID
	ResourceID      string
	Resource        string
	Audience        []string
}

func ParseUploadFilesClaims(tokenStr string, sk string) (UploadFilesJwtData, error) {
	claims := UploadFilesClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(sk), nil
	})
	if err != nil || !token.Valid {
		return UploadFilesJwtData{}, err
	}

	ownerID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return UploadFilesJwtData{}, err
	}

	return UploadFilesJwtData{
		OwnerAccountID:  ownerID,
		UploadSessionID: claims.UploadSessionID,
		ResourceID:      claims.ResourceID,
		Resource:        claims.Resource,
		Audience:        claims.Audience,
	}, nil
}
