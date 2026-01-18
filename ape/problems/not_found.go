package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func NotFound(details string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusNotFound),
		Status: fmt.Sprintf("%d", http.StatusNotFound),
		Code:   "NOT_FOUND",
		Detail: details,
		Meta: &map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	}
}
