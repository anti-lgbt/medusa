package entities

import "time"

type Reply struct {
	ID        int64     `json:"id,omitempty"`
	Content   string    `json:"content,omitempty"`
	LikeCount int64     `json:"like_count,omitempty"`
	Liked     bool      `json:"liked,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
