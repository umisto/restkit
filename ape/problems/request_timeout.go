package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func RequestTimeout() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusRequestTimeout),
		Status: fmt.Sprintf("%d", http.StatusRequestTimeout),
		Code:   "REQUEST_TIMEOUT",
		Meta: &map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	}
}
