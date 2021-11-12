package entities

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	LikeCount int64     `json:"like_count"`
	Replies   []*Reply  `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
