package models

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
)

type Label struct {
	ID        int64            `json:"id" gorm:"primaryKey"`
	UserID    int64            `json:"-" gorm:"type:bigint;not null;uniqueIndex:idx_user_id_and_type;index"`
	Type      string           `json:"type" gorm:"type:character varying(10);uniqueIndex:idx_user_id_and_type;index"`
	State     types.LabelState `json:"state" gorm:"type:character varying(10);not null"`
	CreatedAt time.Time        `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time        `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	User      *User
}
