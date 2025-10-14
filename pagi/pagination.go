package pagi

import (
	"net/http"
	"strconv"
)

func GetPagination(r *http.Request) (page, size uint64) {
	page = 1
	size = 20

	params := r.URL.Query()

	pageStr := params.Get("page")
	if pageStr != "" {
		n, err := strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			page = n
		}
	}

	sizeStr := params.Get("size")
	if sizeStr != "" {
		n, err := strconv.ParseUint(sizeStr, 10, 64)
		if err == nil {
			size = n
		}
	}

	return page, size
}

func PagConvert(page, size uint64) (limit, offset uint64) {
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
