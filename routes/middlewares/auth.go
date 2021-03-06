package middlewares

import (
	"encoding/base64"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
)

func MustAuth(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	user, err := IsAuth(c)
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: err.Error(),
		})
	}

	jwt_token, err := helpers.GenerateJWT(user)
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	session.Set("jwt", jwt_token)
	c.Locals("CurrentUser", user)

	return c.Next()
}

func IsAuth(c *fiber.Ctx) (*models.User, error) {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return nil, errors.New(types.AuthzInvalidSession)
	}

	jwt_token := session.Get("jwt")
	if jwt_token == nil {
		return nil, errors.New(types.AuthzInvalidSession)
	}

	token, err := jwt.Parse(jwt_token.(string), func(t *jwt.Token) (interface{}, error) {
		jwt_private_key_base64 := os.Getenv("JWT_PRIVATE_KEY")

		jwt_private_key, err := base64.StdEncoding.DecodeString(jwt_private_key_base64)
		if err != nil {
			return "", err
		}

		return jwt_private_key, nil
	})
	if err != nil {
		return nil, errors.New(types.JWTDecodeAndVerify)
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var user *models.User
	if result := config.Database.First(&user, "email = ?", email); result.Error != nil {
		session.Destroy()
		return nil, errors.New(types.AuthzInvalidSession)
	}

	return user, nil
}
