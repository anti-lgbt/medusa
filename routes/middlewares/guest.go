package middlewares

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

func IsGuest(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	jwt := session.Get("jwt")
	if jwt != nil {
		return c.Status(422).JSON(types.Error{
			Error: "authz.guest_only",
		})
	}

	return c.Next()
}
