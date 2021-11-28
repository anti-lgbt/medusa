package resource

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/volatiletech/null"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

const (
	UserAvatarNotValid          = "resource.user.avatar_not_valid"
	PasswordDoesntNotMatch      = "resource.password.not_match"
	PasswordPrevPassNotConrrect = "resource.password.prev_pass.not_conrrect"
	PasswordNoChangeProvided    = "resource.password.no_change_provided"
)

// PUT /api/v2/resource/me
func GetUserProfile(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	return c.Status(200).JSON(user.ToEntity())
}

type UpdateUserPayload struct {
	FirstName string      `json:"first_name" form:"first_name" validate:"required"`
	LastName  string      `json:"last_name" form:"last_name" validate:"required"`
	Bio       null.String `json:"bio" form:"bio"`
}

// PUT /api/v2/resource/users
func UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)
	params := new(UpdateUserPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "resource.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	user.FirstName = params.FirstName
	user.LastName = params.LastName
	user.Bio = sql.NullString{
		String: params.Bio.String,
		Valid:  params.Bio.Valid,
	}

	file_header, err := c.FormFile("avatar")
	if err == nil {
		if !services.VerifyFileType(file_header, types.FileTypeImage) {
			return c.Status(422).JSON(types.Error{
				Error: UserAvatarNotValid,
			})
		}

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		user.Avatar = sql.NullString{
			Valid:  true,
			String: file_path,
		}
	}

	config.Database.Save(&user)

	return c.Status(201).JSON(201)
}

type ChangePasswordPayload struct {
	OldPassword     string `json:"old_password" form:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" form:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required"`
}

// PUT /api/v2/resouce/users/password
func UpdateUserPassword(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	params := new(ChangePasswordPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "resource.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	if params.NewPassword != params.ConfirmPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordDoesntNotMatch,
		})
	}

	if services.DecryptPassword(user.Password) != params.OldPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordPrevPassNotConrrect,
		})
	}

	if params.OldPassword == params.NewPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordNoChangeProvided,
		})
	}

	user.UpdatePassword(params.NewPassword)

	services.SendEmail("password_changed", user.Email, user.Language(), nil)

	return c.Status(201).JSON(201)
}
