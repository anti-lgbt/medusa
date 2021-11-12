package entities

import "time"

type Reply struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	LikeCount int64     `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
