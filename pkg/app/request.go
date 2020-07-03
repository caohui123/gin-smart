package app

const (
	URLParamPage     = "page"
	URLParamPageSize = "page_size"
	MaxPageSize      = 200
	DefaultPageSize  = 20
)

type paramsCheck interface {
	Check() error
}

// 有用
type PagerIF interface {
	SetPager(page uint, pageSize, total uint64)
	GetTotal() uint
}

type pagerResp struct {
	Page     uint   `json:"page"`
	PageSize uint   `json:"page_size"`
	Total    uint64 `json:"total"`
}

type Pager struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

type PageResponse struct {
	Pager pagerResp `json:"pager"`
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
	if p.PageSize <= 0 || p.PageSize > MaxPageSize {
		p.PageSize = DefaultPageSize
	}
	return p
}

func (r *PageResponse) GetTotal() uint64 {
	return r.Pager.GetTotal()
}

func (r *PageResponse) SetPager(page uint, pageSize uint, total uint64) {
	r.Pager.SetPager(page, pageSize, total)
}

func (r *pagerResp) SetPager(page uint, pageSize uint, total uint64) {
	r.Page = page
	r.PageSize = pageSize
	r.Total = total
}

func (r *pagerResp) GetTotal() uint64 {
	return r.Total
}

func NewPageRequest(page, pageSize uint) Pager {
	pager := Pager{
		Page:     page,
		PageSize: pageSize,
	}
	pager.secure()
	return pager
}
