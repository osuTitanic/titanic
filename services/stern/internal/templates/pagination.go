package templates

import (
	"net/url"
	"strconv"
)

type PaginationOptions struct {
	Path        string
	Query       url.Values
	CurrentPage int
	PageSize    int
	Total       int
	Window      int
}

type PaginationView struct {
	Total       int
	PageSize    int
	CurrentPage int
	TotalPages  int
	From        int
	To          int
	Items       []PaginationItem
	PreviousUrl string
	NextUrl     string
	HasPrevious bool
	HasNext     bool
}

type PaginationItem struct {
	Label    string
	Page     int
	Url      string
	Current  bool
	Ellipsis bool
}

func NewPagination(options PaginationOptions) PaginationView {
	if options.Path == "" {
		options.Path = "/"
	}
	if options.CurrentPage < 1 {
		options.CurrentPage = 1
	}
	if options.PageSize < 1 {
		options.PageSize = 1
	}

	totalPages := 0
	if options.Total > 0 {
		totalPages = (options.Total + options.PageSize - 1) / options.PageSize
	}
	if totalPages > 0 && options.CurrentPage > totalPages {
		options.CurrentPage = totalPages
	}

	view := PaginationView{
		Total:       options.Total,
		PageSize:    options.PageSize,
		CurrentPage: options.CurrentPage,
		TotalPages:  totalPages,
	}
	if options.Total <= 0 {
		return view
	}

	view.From = (options.CurrentPage-1)*options.PageSize + 1
	view.To = min(options.CurrentPage*options.PageSize, options.Total)
	view.HasPrevious = options.CurrentPage > 1
	view.HasNext = options.CurrentPage < totalPages

	if view.HasPrevious {
		view.PreviousUrl = paginationUrl(options.Path, options.Query, options.CurrentPage-1)
	}
	if view.HasNext {
		view.NextUrl = paginationUrl(options.Path, options.Query, options.CurrentPage+1)
	}

	for _, page := range paginationPages(options.CurrentPage, totalPages) {
		if page == 0 {
			view.Items = append(view.Items, PaginationItem{Ellipsis: true})
			continue
		}
		view.Items = append(view.Items, PaginationItem{
			Label:   strconv.Itoa(page),
			Page:    page,
			Url:     paginationUrl(options.Path, options.Query, page),
			Current: page == options.CurrentPage,
		})
	}
	return view
}

func paginationPages(currentPage, totalPages int) []int {
	if totalPages <= 0 {
		return nil
	}
	if totalPages <= 11 {
		return intRange(1, totalPages)
	}
	if currentPage <= 6 {
		pages := intRange(1, 10)
		return append(pages, 0, totalPages)
	}
	if currentPage >= totalPages-5 {
		pages := []int{1, 0}
		return append(pages, intRange(totalPages-9, totalPages)...)
	}
	pages := []int{1, 0}
	pages = append(pages, intRange(currentPage-4, currentPage+4)...)
	return append(pages, 0, totalPages)
}

func intRange(start, end int) []int {
	values := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		values = append(values, i)
	}
	return values
}

func paginationUrl(path string, current url.Values, page int) string {
	query := cloneQuery(current)
	if page <= 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(page))
	}
	if encoded := query.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

func cloneQuery(current url.Values) url.Values {
	query := make(url.Values, len(current))
	for key, values := range current {
		query[key] = append([]string(nil), values...)
	}
	return query
}
