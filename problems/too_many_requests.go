package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func TooManyRequests() error {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusTooManyRequests),
		Status: fmt.Sprintf("%d", http.StatusTooManyRequests),
		Code:   "TOO_MANY_REQUESTS",
		Meta: &map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	}
}
