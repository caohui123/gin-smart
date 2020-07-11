package util

// 有用
type PagerIF interface {
	SetPager(page uint, pageSize, total uint64)
	GetTotal() uint
}

type PageResp struct {
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
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if  p.PageSize > 500 {
		p.PageSize = 500
	}
	return p
}

func (r *PageResp) SetPager(page uint, pageSize uint, total uint) {
	r.Page = page
	r.PageSize = pageSize
	r.Total = total
}

func (r *PageResp) GetTotal() uint {
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
