package models

import "time"

type Message struct {
	ID        int64
	UserID    int64
	GroupID   int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
