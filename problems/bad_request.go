package problems

import (
	"fmt"
	"io"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
)

type BadRequester interface {
	BadRequest() map[string]error
}

func BadRequest(err error) []error {
	cause := errors.Cause(err)

	if cause == io.EOF {
		return []error{
			&jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Code:   "BAD_REQUEST",
				Detail: "Request body were expected",
				Meta: &map[string]any{
					"timestamp": time.Now().UTC(),
				},
			},
		}
	}

	switch c := cause.(type) {
	case validation.Errors:
		return toJsonapiErrors(c)
	case BadRequester:
		return toJsonapiErrors(c.BadRequest())
	default:
		return []error{
			&jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Code:   "BAD_REQUEST",
				Detail: "Your request was invalid in some way",
				Meta: &map[string]any{
					"timestamp": time.Now().UTC(),
				},
			},
		}
	}
}

func toJsonapiErrors(m map[string]error) []error {
	errs := make([]error, 0, len(m))

	for key, value := range m {
		if value == nil {
			continue
		}

		errs = append(errs, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Code:   "BAD_REQUEST",
			Meta: &map[string]any{
				"field":     key,
				"error":     value.Error(),
				"timestamp": time.Now().UTC(),
			},
		})
	}

	return errs
}
