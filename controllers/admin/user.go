package admin

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/admin/entities"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/controllers/queries"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/creasty/defaults"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/volatiletech/null"
)

const (
	UserAvatarNotValid = "admin.user.avatar_not_valid"
)

func UserToEntity(user *models.User) *entities.User {
	return &entities.User{
		ID:        user.ID,
		UID:       user.UID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bio: null.String{
			String: user.Bio.String,
			Valid:  user.Bio.Valid,
		},
		State:     user.State,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type GetUserPayload struct {
	State types.UserState `json:"state" validate:"userState"`
	Role  types.UserRole  `json:"role" validate:"userRole"`
	queries.Pagination
	queries.Period
	queries.Order
}

func GetUsers(c *fiber.Ctx) error {
	var users []*models.User

	params := new(GetUserPayload)
	if c.QueryParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	defaults.Set(params)

	if err := helpers.Vaildate(params, "admin.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	tx := config.Database
	tx = queries.QueryPagination(tx, params.Limit, params.Page)
	tx = queries.QueryOrder(tx, params.OrderBy, params.Ordering)
	tx = queries.QueryPeriod(tx, params.TimeFrom, params.TimeTo)

	if len(params.State) > 0 {
		tx = tx.Where("state = ?", params.State)
	}

	if len(params.Role) > 0 {
		tx = tx.Where("role = ?", params.Role)
	}

	tx.Find(&users)

	user_entities := make([]*entities.User, 0)

	for _, user := range users {
		user_entities = append(user_entities, UserToEntity(user))
	}

	return c.Status(200).JSON(user_entities)
}

// Get /api/v2/admin/users/:uid
func GetUser(c *fiber.Ctx) error {
	uid := c.Params("uid")

	var user *models.User
	if result := config.Database.First(&user, "uid = ?", uid); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(201).JSON(UserToEntity(user))
}

type UpdateUserPayload struct {
	FirstName string          `json:"first_name" form:"first_name" validate:"required"`
	LastName  string          `json:"last_name" form:"last_name" validate:"required"`
	Bio       null.String     `json:"bio" form:"bio"`
	State     types.UserState `json:"state" form:"state" validate:"userState|required"`
	Role      types.UserRole  `json:"role" form:"role" validate:"userRole|required"`
}

// PUT /api/v2/admin/users/:uid
func UpdateUser(c *fiber.Ctx) error {
	params := new(UpdateUserPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "admin.user"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	var user *models.User
	uid := c.Params("uid")
	if result := config.Database.First(&user, "uid = ?", uid); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	user.FirstName = params.FirstName
	user.LastName = params.LastName
	user.Bio = sql.NullString{
		String: params.Bio.String,
		Valid:  params.Bio.Valid,
	}
	user.State = params.State
	user.Role = params.Role

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

	return c.Status(201).JSON(UserToEntity(user))
}

// DELETE /api/v2/admin/users/:uid
func DeleteUser(c *fiber.Ctx) error {
	var user *models.User
	uid := c.Params("uid")
	if result := config.Database.First(&user, "uid = ?", uid); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	config.Database.Delete(&user)

	return c.Status(200).JSON(200)
}
