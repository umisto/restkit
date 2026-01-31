package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

func InternalError() error {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusInternalServerError),
		Status: fmt.Sprintf("%d", http.StatusInternalServerError),
		Code:   "INTERNAL_SERVER_ERROR",
		Detail: "Oops, an unexpected error has occurred." +
			" We are already looking into it and doing everything we can to ensure that you don't see it again." +
			" Please come back soon.",
		Meta: &map[string]any{
			"timestamp": time.Now().UTC(),
		},
	}
}
