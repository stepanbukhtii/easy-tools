package api

var DefaultPerPage = 10

type Pagination struct {
	Page    int `form:"page" json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}

func (p Pagination) Limit() int {
	if p.PerPage > 0 {
		return p.PerPage
	}
	return DefaultPerPage
}

func (p Pagination) Offset() int {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PerPage > 0 {
		return p.PerPage * (p.Page - 1)
	}

	return DefaultPerPage * (p.Page - 1)
}

func (p Pagination) LimitOffset() (int, int) {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PerPage > 0 {
		return p.PerPage, p.PerPage * (p.Page - 1)
	}

	return DefaultPerPage, DefaultPerPage * (p.Page - 1)
}

func (p Pagination) Pages(total int64) int64 {
	if p.PerPage > 0 {
		return ((total - 1) / int64(p.PerPage)) + 1
	}
	return ((total - 1) / int64(DefaultPerPage)) + 1
}
