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
	UserEmailAlreadyExist    = "identity.user.already_exist"
	UserFailedToCreate       = "identity.user.failed_to_create"
	UserFailedToSendingEmail = "identity.user.failed_to_sending_email"
	UserCodeInconrrect       = "identity.user.code_incorrect"
)

type RegisterPayload struct {
	FirstName string `json:"first_name" form:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" validate:"required"`
	LoginPayload
}

func Register(c *fiber.Ctx) error {
	params := new(RegisterPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "identity.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error == nil {
		return c.Status(422).JSON(types.Error{
			Error: UserEmailAlreadyExist,
		})
	}

	uid_number, err := models.GenerateCode(10)
	if err != nil {
		return c.Status(422).JSON(types.Error{
			Error: UserFailedToCreate,
		})
	}

	user = &models.User{
		UID:       "UID" + uid_number,
		Email:     params.Email,
		Password:  services.EncryptPassword(params.Password),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		State:     types.UserStatePending,
		Role:      types.UserRoleMember,
	}

	config.Database.Create(&user)

	label := models.Label{
		UserID: user.ID,
		Type:   "email",
		State:  types.LabelStatePending,
	}

	config.Database.Create(&label)

	code := user.GetConfirmationCode("email", true)

	code.SendCode("email_confirmation", user.Language())

	session, err := config.SessionStore.Get(c)
	if err != nil {
		return c.Status(422).JSON(types.Error{
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
	session.Save()

	return c.Status(201).JSON(user.ToEntity())
}

type ResendEmailCode struct {
	Email string `json:"email" form:"email"`
}

func ReSendEmailCode(c *fiber.Ctx) error {
	params := new(ResendEmailCode)
	if err := c.BodyParser(params); err != nil {
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

	code := user.GetConfirmationCode("email", true)

	code.SendCode("email_confirmation", user.Language())

	return c.Status(200).JSON(200)
}

type VerifyEmailPayload struct {
	Email string `json:"email" form:"email" validate:"required|email"`
	Code  string `json:"code" form:"code" validate:"required"`
}

func VerifyEmail(c *fiber.Ctx) error {
	params := new(VerifyEmailPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "identity.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	var user *models.User
	if result := config.Database.First(&user, "email = ?", params.Email); result.Error != nil {
		return c.Status(422).JSON(types.Error{
			Error: UserInvalid,
		})
	}

	code := user.GetConfirmationCode("email", false)

	if code.Expired() || code.Validated() || code.OutAttempt() {
		return c.Status(422).JSON(types.Error{
			Error: UserCodeInconrrect,
		})
	}

	if code.Code != params.Code {
		code.AttemptCount++

		config.Database.Save(&code)
		return c.Status(422).JSON(types.Error{
			Error: UserCodeInconrrect,
		})
	}

	code.Validation()

	user.State = types.UserStateActive

	// for _, label := range user.Labels {
	// 	label.State = types.LabelStateVerified
	// }

	config.Database.Save(&user)

	code.SendCode("email_verification_successful", user.Language())

	session, err := config.SessionStore.Get(c)
	if err != nil {
		return c.Status(422).JSON(types.Error{
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
	session.Save()

	return c.Status(200).JSON(200)
}
