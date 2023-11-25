package api

var DefaultPerPage int64 = 10

type Pagination struct {
	Page    int64 `form:"page" json:"page"`
	PerPage int64 `form:"per_page" json:"per_page"`
}

func (p Pagination) LimitOffset() (int64, int64) {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PerPage > 0 {
		return p.PerPage, p.PerPage * (p.Page - 1)
	}

	return DefaultPerPage, DefaultPerPage * (p.Page - 1)
}

func (p Pagination) Pages(count int64) int64 {
	if p.PerPage > 0 {
		return ((count - 1) / p.PerPage) + 1
	}
	return ((count - 1) / DefaultPerPage) + 1
}
