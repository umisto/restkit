package pagi

import (
	"net/http"
	"strings"
)

type SortField struct {
	Field  string
	Ascend bool // "asc" or "desc"
}

func SortFields(r *http.Request) (sortFields []SortField) {
	if sortStr := r.URL.Query().Get("sort"); sortStr != "" {
		parts := strings.Split(sortStr, ",")
		for _, p := range parts {
			ascend := true
			field := p
			if strings.HasPrefix(p, "-") {
				ascend = false
				field = strings.TrimPrefix(p, "-")
			}
			sortFields = append(sortFields, SortField{
				Field:  field,
				Ascend: ascend,
			})
		}
	}

	return sortFields
}
