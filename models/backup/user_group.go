package models

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
)

type UserGroup struct {
	ID        int64
	UserID    int64
	GroupID   int64
	Role      types.GroupRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
