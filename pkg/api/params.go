package api

import (
	"fmt"
	"strconv"
)

type PathParam struct {
	CmdArgs []string
}

func (p *PathParam) GetInt64(idx int) (int64, error) {
	if idx >= len(p.CmdArgs) {
		return 0, fmt.Errorf("idx %d exceed", idx)
	}

	return strconv.ParseInt(p.CmdArgs[idx], 10, 64)
}

func (p *PathParam) GetString(idx int) (string, error) {
	if idx >= len(p.CmdArgs) {
		return "", fmt.Errorf("idx %d exceed", idx)
	}

	return p.CmdArgs[idx], nil
}

const (
	DefaultPageSize = 20
	MaxPageSize     = 200
)

type PageParams struct {
	Page      int64  `json:"page" form:"page" validate:"gt=0" reg_error_info:"参数有误"`           // 页码
	PageSize  int64  `json:"page_size" form:"page_size" validate:"gt=0" reg_error_info:"参数有误"` // 分页大小
	PageToken string `json:"page_token" form:"page_token"`                                     // 分页标识，PageToken 时，不使用page 参数
	Sort      string `form:"sort" validate:"oneof=desc asc ''" reg_error_info:"只能选desc/asc"`
}

type PageResult struct {
	TotalCount int64       `json:"total_count"` // 记录总数
	Page       int64       `json:"page"`        // 当前页码，当使用 PageCursor 时，不使用
	PageSize   int64       `json:"page_size"`   // 分页大小
	PageToken  string      `json:"page_token"`  // 分页token，下一页时传递该游标，推荐使用
	List       interface{} `json:"list"`
}

func WrapPageResult(params *PageParams, data any, totalCount int64) *PageResult {
	return &PageResult{
		TotalCount: totalCount,
		List:       data,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}
}

func (p *PageParams) Offset() int64 {
	p.normal()

	return p.PageSize * (p.Page - 1)
}

func (p *PageParams) Limit() int64 {
	p.normal()

	return p.PageSize
}

func (p *PageParams) normal() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}

	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}
