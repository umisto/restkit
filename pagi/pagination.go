package pagi

import (
	"net/http"
	"strconv"
)

func GetPagination(r *http.Request) (page, size int32) {
	page = 1
	size = 20

	params := r.URL.Query()

	pageStr := params.Get("page")
	if pageStr != "" {
		n, err := strconv.ParseInt(pageStr, 10, 32)
		if err == nil {
			page = int32(n)
		}
	}

	sizeStr := params.Get("size")
	if sizeStr != "" {
		n, err := strconv.ParseInt(sizeStr, 10, 32)
		if err == nil {
			size = int32(n)
		}
	}

	return page, size
}

func PagConvert(page, size int32) (limit, offset int32) {
	if page == 0 {
		page = 1
	}
	if size > 100 {
		size = 100
	}
	if size == 0 {
		size = 20
	}

	limit = size
	offset = (page - 1) * size
	return
}
