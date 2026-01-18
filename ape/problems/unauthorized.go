package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func Unauthorized(details string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusUnauthorized),
		Status: fmt.Sprintf("%d", http.StatusUnauthorized),
		Code:   "UNAUTHORIZED",
		Detail: details,
		Meta: &map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	}
}
