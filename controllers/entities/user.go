package entities

import (
	"time"

	"github.com/anti-lgbt/medusa/types"
	"github.com/volatiletech/null"
)

type User struct {
	ID        int64           `json:"id"`
	UID       string          `json:"uid"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Bio       null.String     `json:"bio"`
	State     types.UserState `json:"state"`
	Role      types.UserRole  `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
