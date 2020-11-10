package models

import (
	xtime "gin-template-app/pkg/time"
	"time"
)

// Page page query type
type Page struct {
	Page  *uint64 `form:"page" json:"page"`   // 页数
	Limit *uint64 `form:"limit" json:"limit"` // 页长
	Sort  *string `form:"sort" json:"sort"`   // 排序
}

// Detail KV
type Detail struct {
	ID     uint   `json:"id" form:"id"`
	Common string `json:"common" form:"common"`
}

// Base Basic database data
type Base struct {
	ID        uint      `gorm:"type:serial;primary_key;auto_increment;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BaseView basic database data view
type BaseView struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NumberQuery Fuzzy query number type
type NumberQuery struct {
	Max int `form:"max" json:"max"`
	Min int `form:"min" json:"min"`
}

// PageQuery page query type
type PageQuery struct {
	Page  *uint64 `form:"page" json:"page"`
	Limit *uint64 `form:"limit" json:"limit"`
	Sort  *string `form:"sort" json:"sort"`
}

type timeQuery struct {
	Start xtime.Time `form:"start" json:"start"`
	End   xtime.Time `form:"end" json:"end"`
}

// TimeQuery time query
type TimeQuery struct {
	CreatedAt *timeQuery `form:"created_at" json:"created_at"`
	UpdateAt  *timeQuery `form:"update_at" json:"update_at"`
}

// ValidPageForm valid page query form
func ValidPageForm(pageForm *Page) {
	// PageQuery to PageWhereOrder
	if pageForm.Page == nil || *pageForm.Page < 1 {
		pageForm.Page = new(uint64)
		*pageForm.Page = 1
	}
	if pageForm.Limit == nil {
		pageForm.Limit = new(uint64)
		*pageForm.Limit = 10
	}
	if pageForm.Sort != nil {
		sort := *pageForm.Sort
		if len(sort) > 2 {
			orderType := sort[0:1]
			order := sort[1:]
			if orderType == "^" {
				order += " ASC"
			} else {
				order += " DESC"
			}
			*pageForm.Sort = order
		} else {
			pageForm.Sort = new(string)
			*pageForm.Sort = "id ASC"
		}
	} else {
		pageForm.Sort = new(string)
		*pageForm.Sort = "id ASC"
	}
}
