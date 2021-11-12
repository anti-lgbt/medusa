package entities

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
	"github.com/volatiletech/null"
)

type Music struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description null.String      `json:"description"`
	State       types.MusicState `json:"state"`
	ViewCount   int64            `json:"view_count"`
	LikeCount   int64            `json:"like_count"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
