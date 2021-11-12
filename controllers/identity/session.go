package identity

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

const (
	UserInvalid = "identity.session.invalid_params"
	UserDeleted = "identity.user.deleted"
	UserBanned  = "identity.user.banned"
)

type LoginPayload struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func Login(c *fiber.Ctx) error {
	var params *LoginPayload
	if err := c.BodyParser(&params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error != nil {
		return c.Status(422).JSON(types.Error{
			Error: UserInvalid,
		})
	}

	if user.State == types.UserStateDeleted {
		return c.Status(401).JSON(types.Error{
			Error: UserDeleted,
		})
	}

	if user.State == types.UserStateBanned {
		return c.Status(401).JSON(types.Error{
			Error: UserBanned,
		})
	}

	if user.DecryptedPassword() != params.Password {
		return c.Status(422).JSON(types.Error{
			Error: UserDeleted,
		})
	}

	session, err := config.SessionStore.Get(c)
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: UserInvalid,
		})
	}

	jwt_token, err := helpers.GenerateJWT(user)
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	session.Set("jwt", jwt_token)

	services.SendEmail("email_verification_successful", user.Language(), nil)

	return c.Status(200).JSON(user.ToEntity())
}
