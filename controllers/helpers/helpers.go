package helpers

import (
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

	return token.SignedString(os.Getenv("JWT_PRIVATE_KEY"))
}

func VerifyFileType(file_header *multipart.FileHeader, file_type string) bool {
	return true
}
