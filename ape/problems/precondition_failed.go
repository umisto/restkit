package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func PreconditionFailed(details string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusPreconditionFailed),
		Status: fmt.Sprintf("%d", http.StatusPreconditionFailed),
		Code:   "STATUS_PRECONDITION_FAILED",
		Detail: details,
		Meta: &map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	}
}
