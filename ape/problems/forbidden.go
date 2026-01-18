package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func Forbidden(details string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusForbidden),
		Status: fmt.Sprintf("%d", http.StatusForbidden),
		Code:   "FORBIDDEN",
		Detail: details,
		Meta: &map[string]any{
			"timestamp": time.Now().UTC(),
		},
	}
}
