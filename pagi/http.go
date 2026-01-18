package pagi

import (
	"net/http"
	"strconv"
	"strings"
)

func GetLimit(
	r *http.Request,
	max uint,
) uint {
	limit := uint(20)

	if s := r.URL.Query().Get("limit"); s != "" {
		parsed, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			limit = uint(parsed)
		}
	}

	if limit == 0 {
		limit = 20
	}
	if limit > max {
		limit = max
	}

	return limit
}

func GetSort(r *http.Request) *SortField {
	sortStr := strings.TrimSpace(r.URL.Query().Get("sort"))
	if sortStr == "" {
		return nil
	}

	sort := &SortField{
		Ascend: true,
	}

	if strings.HasPrefix(sortStr, "-") {
		sort.Ascend = false
		sort.Field = strings.TrimPrefix(sortStr, "-")
	} else {
		sort.Field = sortStr
	}

	sort.Field = strings.TrimSpace(sort.Field)
	if sort.Field == "" {
		return nil
	}

	return sort
}

func GetPagination(r *http.Request) (limit uint, offset uint) {
	page := uint(1)
	size := uint(20)

	params := r.URL.Query()

	if s := params.Get("page"); s != "" {
		if n, err := strconv.ParseUint(s, 10, 32); err == nil && n > 0 {
			page = uint(n)
		}
	}

	if s := params.Get("size"); s != "" {
		if n, err := strconv.ParseUint(s, 10, 32); err == nil && n > 0 {
			size = uint(n)
		}
	}

	const maxSize = uint(100)
	if size > maxSize {
		size = maxSize
	}

	offset = (page - 1) * size
	limit = size

	return
}

type Links struct {
	Self  string  `json:"self"`
	First *string `json:"first,omitempty"`
	Last  *string `json:"last,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Next  *string `json:"next,omitempty"`
}

func BuildPageLinks(r *http.Request, page, size, total uint) Links {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}

	lastPage := uint(1)
	if total > 0 {
		lastPage = (total + size - 1) / size
		if lastPage == 0 {
			lastPage = 1
		}
	}

	self := buildURLWithPage(r, page, size)

	var first *string
	if page != 1 {
		v := buildURLWithPage(r, 1, size)
		first = &v
	}

	var last *string
	if page != lastPage {
		v := buildURLWithPage(r, lastPage, size)
		last = &v
	}

	var prev *string
	if page > 1 {
		v := buildURLWithPage(r, page-1, size)
		prev = &v
	}

	var next *string
	if page < lastPage {
		v := buildURLWithPage(r, page+1, size)
		next = &v
	}

	return Links{
		Self:  self,
		First: first,
		Last:  last,
		Prev:  prev,
		Next:  next,
	}
}

func buildURLWithPage(r *http.Request, page, size uint) string {
	u := *r.URL
	q := u.Query()

	q.Set("page", strconv.FormatUint(uint64(page), 10))
	q.Set("size", strconv.FormatUint(uint64(size), 10))

	q.Del("page[number]")
	q.Del("page[size]")

	u.RawQuery = q.Encode()
	return u.String()
}
