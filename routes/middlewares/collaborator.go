package middlewares

import (
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

func MustCollaborator(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	for _, role := range []types.UserRole{types.UserRoleSinger, types.UserRoleMusician, types.UserRoleAdmin} {
		if user.Role == role {
			return c.Next()
		}
	}

	return c.Status(401).JSON(types.Error{
		Error: types.AuthZInvalidPermission,
	})
}
