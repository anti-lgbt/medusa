package entities

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
	"github.com/volatiletech/null"
)

type Music struct {
	ID          int64            `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Author      string           `json:"author,omitempty"`
	Description null.String      `json:"description,omitempty"`
	State       types.MusicState `json:"state,omitempty"`
	ViewCount   int64            `json:"view_count,omitempty"`
	LikeCount   int64            `json:"like_count,omitempty"`
	Liked       bool             `json:"liked,omitempty"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
}
