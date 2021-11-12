package entities

import (
	"time"

	"github.com/anti-lgbt/medusa/models/datatypes"
)

type Album struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Description datatypes.NullString `json:"description"`
	Private     bool                 `json:"private"`
	ViewCount   int64                `json:"view_count"`
	LikeCount   int64                `json:"like_count"`
	Musics      []*Music             `json:"music"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}
