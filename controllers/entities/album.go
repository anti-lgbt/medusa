package entities

import (
	"time"

	"github.com/volatiletech/null"
)

type Album struct {
	ID          int64       `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description null.String `json:"description,omitempty"`
	Private     bool        `json:"private,omitempty"`
	ViewCount   int64       `json:"view_count,omitempty"`
	LikeCount   int64       `json:"like_count,omitempty"`
	Musics      []*Music    `json:"music,omitempty"`
	Liked       bool        `json:"liked,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
}
