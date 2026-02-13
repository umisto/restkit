package headers

import (
	"fmt"
	"net/http"
	"strings"
)

const UploadContentHeader = "Upload-Token"

func GetUploadContent(r *http.Request) (string, error) {
	authHeader := r.Header.Get(UploadContentHeader)
	if authHeader == "" {
		return "", fmt.Errorf("missing %s header", UploadContentHeader)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid %s header format", UploadContentHeader)
	}

	return parts[1], nil
}
