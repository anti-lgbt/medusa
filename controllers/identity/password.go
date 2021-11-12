package identity

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

const (
	PasswordCodeInconrrect = "identity.password.code_incorrect"
	PasswordDoesntNotMatch = "identity.password.does_not_match"
)

func GenerateCodeResetPassword(c *fiber.Ctx) error {
	type Payload struct {
		Email string `json:"email" form:"email"`
	}

	params := new(Payload)
	if c.BodyParser(&params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	code := user.GetConfirmationCode("email", true)

	code.SendCode("password_reset", user.Language())

	return c.Status(200).JSON(200)
}

func CheckCodeResetPassword(c *fiber.Ctx) error {
	type Payload struct {
		Email string `json:"email" form:"email"`
		Code  string `json:"code" form:"code"`
	}

	params := new(Payload)
	if c.BodyParser(&params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	code := user.GetConfirmationCode("email", true)

	if code.Expired() || code.Validated() || code.OutAttempt() {
		return c.Status(422).JSON(types.Error{
			Error: PasswordCodeInconrrect,
		})
	}

	if code.Code != params.Code {
		code.AttemptCount++

		return c.Status(422).JSON(types.Error{
			Error: PasswordCodeInconrrect,
		})
	}

	return c.Status(200).JSON(200)
}

func ResetPassword(c *fiber.Ctx) error {
	type Payload struct {
		Email           string `json:"email" form:"email"`
		Code            string `json:"code" form:"code"`
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	params := new(Payload)
	if c.BodyParser(&params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	code := user.GetConfirmationCode("email", true)

	if code.Expired() || code.Validated() || code.OutAttempt() {
		return c.Status(422).JSON(types.Error{
			Error: PasswordCodeInconrrect,
		})
	}

	if code.Code != params.Code {
		code.AttemptCount++

		return c.Status(422).JSON(types.Error{
			Error: PasswordCodeInconrrect,
		})
	}

	if params.Password != params.ConfirmPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordDoesntNotMatch,
		})
	}

	user.UpdatePassword(params.Password)
	code.Validation()

	services.SendEmail("password_reset_successful", user.Email, user.Language(), nil)

	return c.Status(200).JSON(200)
}
