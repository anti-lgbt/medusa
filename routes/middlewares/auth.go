package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
)

func IsAuth(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	jwt_token := session.Get("jwt")
	if jwt_token == nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	token, err := jwt.Parse(jwt_token.(string), func(t *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_PRIVATE_KEY"), nil
	})
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.JWTDecodeAndVerify,
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var user *models.User
	if result := config.Database.First(&user, "email = ?", email); result.Error != nil {
		session.Destroy()
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	jwt_token, err = helpers.GenerateJWT(user)
	if err != nil {
		return c.Status(401).JSON(types.Error{
			Error: types.AuthzInvalidSession,
		})
	}

	session.Set("jwt", jwt_token)
	c.Locals("CurrentUser", user)

	return c.Next()
}
