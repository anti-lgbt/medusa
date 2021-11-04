package middlewares

import (
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

func IsAdmin(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	if user.Role != types.UserRoleAdmin {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthZInvalidPermission,
		})
	}

	return c.Next()
}
