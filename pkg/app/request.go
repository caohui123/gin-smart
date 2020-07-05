package app




// 有用
type PagerIF interface {
	SetPager(page uint, pageSize, total uint64)
	GetTotal() uint
}

type pagerResp struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
	Total    uint `json:"total"`
}

type Pager struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

func (p *Pager) Offset() uint {
	p.secure()
	return (p.Page - 1) * p.PageSize
}

func (p *Pager) Limit() uint {
	p.secure()
	return p.PageSize
}

func (p *Pager) secure() *Pager {
	defaultPageSize := Runner.Cfg.General.DefaultPageSize
	maxPageSize := Runner.Cfg.General.MaxPageSize
	if maxPageSize == 0 {
		maxPageSize = 200
	}
	if defaultPageSize == 0 {
		defaultPageSize = 20
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > maxPageSize {
		p.PageSize = defaultPageSize
	}
	return p
}

func (r *pagerResp) SetPager(page uint, pageSize uint, total uint) {
	r.Page = page
	r.PageSize = pageSize
	r.Total = total
}

func (r *pagerResp) GetTotal() uint {
	return r.Total
}

func NewRequestPager(page, pageSize uint) *Pager {
	pager := &Pager{
		Page:     page,
		PageSize: pageSize,
	}
	pager.secure()
	return pager
}
