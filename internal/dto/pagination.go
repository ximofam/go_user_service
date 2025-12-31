package dto

import "strings"

type PagingOutput struct {
	CurrentPage int   `json:"current_Page"`
	TotalPages  int   `json:"total_Pages"`
	Limit       int   `json:"Limit"`
	Total       int64 `json:"total"`
}

type PagingInput struct {
	Page    int    `form:"page" json:"-"`
	Limit   int    `form:"limit" json:"-"`
	SortBy  string `form:"sort_by" json:"-"`
	SortDir string `form:"sort_dir" json:"-"`
}

func (p *PagingInput) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PagingInput) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

func (p *PagingInput) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PagingInput) GetSortDirection() string {
	dir := strings.ToUpper(p.SortDir)
	if dir == "ASC" {
		return "ASC"
	}
	return "DESC"
}

// func (p *PagingInput) GetSortColumn(allowedCols map[string]bool) string {
// 	col := strings.ToLower(p.SortBy)

// 	if allowedCols[col] {
// 		return col
// 	}
// 	return "id"
// }

func (p *PagingInput) GetSortColumn() string {
	if p.SortBy == "" {
		return "id"
	}

	return strings.ToLower(p.SortBy)
}
