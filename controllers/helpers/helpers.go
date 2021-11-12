package helpers

import (
	"encoding/base64"
	"mime/multipart"
	"os"

	"github.com/anti-lgbt/medusa/models"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":        user.UID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"state":      user.State,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt_private_key_base64 := os.Getenv("JWT_PRIVATE_KEY")

	jwt_private_key, err := base64.StdEncoding.DecodeString(jwt_private_key_base64)
	if err != nil {
		return "", err
	}

	return token.SignedString([]byte(jwt_private_key))
}

func VerifyFileType(file_header *multipart.FileHeader, file_type string) bool {
	return true
}
