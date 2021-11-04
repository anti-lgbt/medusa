package models

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        int64
	UID       string // Random or user to user
	Personal  bool   // True for user to user false for group
	Name      sql.NullString
	Avatar    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
