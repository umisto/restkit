package problems

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
	"github.com/pkg/errors"
)

func isBadRequest(err error) bool {
	e, ok := err.(interface {
		BadRequest() bool
	})
	return ok && e.BadRequest()
}

func isNotAllowed(err error) bool {
	e, ok := err.(interface {
		NotAllowed() bool
	})
	return ok && e.NotAllowed()
}

func isForbidden(err error) bool {
	e, ok := err.(interface {
		Forbidden() bool
	})
	return ok && e.Forbidden()
}

// NotAllowed will try to guess details of error and populate problem accordingly.
func NotAllowed(details string, errs ...error) error {
	// errs is optional for backward compatibility
	if len(errs) == 0 {
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusUnauthorized),
			Status: fmt.Sprintf("%d", http.StatusUnauthorized),
		}
	}

	if len(errs) != 1 {
		panic(errors.New("unexpected number of errors passed"))
	}

	cause := errors.Cause(errs[0])
	switch {
	case isBadRequest(cause):
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Code:   "BAD_REQUEST",
			Detail: details,
			Meta: &map[string]interface{}{
				"reason":    cause.Error(),
				"timestamp": time.Now().UTC(),
			},
		}
	case isNotAllowed(cause):
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusUnauthorized),
			Status: fmt.Sprintf("%d", http.StatusUnauthorized),
			Code:   "NOT_ALLOWED",
			Detail: details,
			Meta: &map[string]interface{}{
				"timestamp": time.Now().UTC(),
			},
		}
	case isForbidden(cause):
		{
			return &jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusForbidden),
				Status: fmt.Sprintf("%d", http.StatusForbidden),
				Code:   "FORBIDDEN",
				Detail: details,
				Meta: &map[string]interface{}{
					"timestamp": time.Now().UTC(),
				},
			}
		}
	default:
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusInternalServerError),
			Status: fmt.Sprintf("%d", http.StatusInternalServerError),
			Code:   "INTERNAL_SERVER_ERROR",
			Detail: details,
			Meta: &map[string]interface{}{
				"timestamp": time.Now().UTC(),
			},
		}
	}
}
