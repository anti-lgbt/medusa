package entities

import (
	"time"

	"github.com/anti-lgbt/medusa/models/datatypes"
	"github.com/anti-lgbt/medusa/types"
)

type User struct {
	ID        int64                `json:"id"`
	UID       string               `json:"uid"`
	Email     string               `json:"email"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Bio       datatypes.NullString `json:"bio"`
	State     types.UserState      `json:"state"`
	Role      types.UserRole       `json:"role"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
