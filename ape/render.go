package ape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"errors"

	"github.com/google/jsonapi"
)

func Render(w http.ResponseWriter, status int, res ...interface{}) {
	w.Header().Set("content-type", jsonapi.MediaType)
	w.WriteHeader(status)

	if res != nil && len(res) > 0 {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			panic(fmt.Errorf("failed to render response %w", err))
		}
	}
}

func RenderErr(w http.ResponseWriter, errs ...*jsonapi.ErrorObject) {
	if len(errs) == 0 {
		panic("expected non-empty errors slice")
	}

	objs := make([]*jsonapi.ErrorObject, 0, len(errs))

	for _, e := range errs {
		if e == nil {
			continue
		}

		var jo *jsonapi.ErrorObject
		if errors.As(e, &jo) && jo != nil {
			objs = append(objs, jo)
			continue
		}
	}

	if len(objs) == 0 {
		panic("no renderable jsonapi errors produced")
	}

	status, err := strconv.Atoi(objs[0].Status)
	if err != nil {
		log.Printf("Failed to parse status: %v", err)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)

	if err = jsonapi.MarshalErrors(w, objs); err != nil {
		log.Printf("Failed to marshal errors: %v", err)
	}
}
