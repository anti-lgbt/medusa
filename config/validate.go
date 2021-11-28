package config

import (
	"github.com/anti-lgbt/medusa/types"
	"github.com/gookit/validate"
)

func InitValidator() {
	validate.AddValidator("userRole", func(role types.UserRole) bool {
		for _, r := range []types.UserRole{types.UserRoleAdmin, types.UserRoleMember, types.UserRoleMusician} {
			if r == role {
				return true
			}
		}

		return false
	})

	validate.AddValidator("userState", func(state types.UserState) bool {
		for _, s := range []types.UserState{types.UserStateActive, types.UserStateBanned, types.UserStateDeleted, types.UserStatePending} {
			if s == state {
				return true
			}
		}

		return false
	})

	t := validate.NewTranslator()
	t.AddFieldMap(validate.MS{
		"State":    "state",
		"Role":     "role",
		"Type":     "type",
		"Limit":    "limit",
		"Page":     "page",
		"TimeFrom": "time_from",
		"TimeTo":   "time_to",
		"OrderBy":  "order_by",
	})
}
