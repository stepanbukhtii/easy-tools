package api

var PaginationPerPage int64 = 10

type Pagination struct {
	Page    int64 `form:"page" json:"page"`
	PerPage int64 `form:"per_page" json:"per_page"`
}

func (p Pagination) LimitOffset() (int64, int64) {
	if p.Page == 0 {
		p.Page = 1
	}

	limit := PaginationPerPage
	if p.PerPage > 0 {
		limit = p.PerPage
	}

	return limit, limit * (p.Page - 1)
}

func (p Pagination) Pages(count int64) int64 {
	limit := PaginationPerPage
	if p.PerPage > 0 {
		limit = p.PerPage
	}
	return ((count - 1) / limit) + 1
}
