package queries

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit int `query:"limit" validate:"uint" default:"100"`
	Page  int `query:"page" validate:"uint" default:"1"`
}

type Period struct {
	TimeFrom int64 `query:"time_from" validate:"uint"`
	TimeTo   int64 `query:"time_to" validate:"uint"`
}

type Order struct {
	OrderBy  string         `query:"order_by" default:"id"`
	Ordering types.Ordering `query:"ordering" default:"asc"`
}

func QueryPagination(tx *gorm.DB, limit, page int) *gorm.DB {
	return tx.Offset(page*limit - limit).Limit(limit)
}

func QueryPeriod(tx *gorm.DB, time_from, time_to int64) *gorm.DB {
	if time_from > 0 {
		tx = tx.Where("created_at >= ?", time.Unix(time_from, 0))
	}

	if time_to > 0 {
		tx = tx.Where("updated_at >= ?", time.Unix(time_to, 0))
	}

	return tx
}

func QueryOrder(tx *gorm.DB, order_by string, ordering types.Ordering) *gorm.DB {
	if len(order_by) > 0 {
		tx = tx.Order(order_by + " " + string(ordering))
	}

	return tx
}
