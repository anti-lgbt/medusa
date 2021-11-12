package models

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/volatiletech/null"
)

type UserData struct {
	Language string   `json:"language"`
	IPS      []string `json:"ips"`
}

type User struct {
	ID         int64           `gorm:"primaryKey;autoIncrement;not null;index"`
	UID        string          `gorm:"type:character varying(20);not null;index"`
	Email      string          `gorm:"type:character varying(50);not null;uniqueIndex"`
	Password   string          `gorm:"type:text;not null"`
	FirstName  string          `gorm:"type:character varying(50);not null;index"`
	LastName   string          `gorm:"type:character varying(50);not null;index"`
	Bio        sql.NullString  `gorm:"type:text"`
	State      types.UserState `gorm:"type:character varying(10);not null;index"`
	Role       types.UserRole  `gorm:"type:character varying(10);not null;index"`
	Avatar     sql.NullString  `gorm:"type:text"`
	Data       sql.NullString  `gorm:"type:text"`
	CreatedAt  time.Time       `gorm:"type:timestamp(0);not null;index"`
	UpdatedAt  time.Time       `gorm:"type:timestamp(0);not null;index"`
	Populars   []*Popular      `gorm:"constraint:OnDelete:CASCADE"`
	Activities []*Activity     `gorm:"constraint:OnDelete:CASCADE"`
	Labels     []*Label        `gorm:"constraint:OnDelete:CASCADE"`
	Musics     []*Music        `gorm:"constraint:OnDelete:CASCADE"`
	Albums     []*Album        `gorm:"constraint:OnDelete:CASCADE"`
	Likes      []*Like         `gorm:"constraint:OnDelete:CASCADE"`
	Comments   []*Comment      `gorm:"constraint:OnDelete:CASCADE"`
	Replys     []*Reply        `gorm:"constraint:OnDelete:CASCADE"`
	Codes      []*Code         `gorm:"constraint:OnDelete:CASCADE"`
}

func (u *User) Language() string {
	var user_data *UserData

	if u.Data.Valid && json.Unmarshal([]byte(u.Data.String), &user_data) != nil {
		return os.Getenv("DEFAULT_LANGUAGE")
	}

	return user_data.Language
}

func (u *User) GetConfirmationCode(code_type string, reset bool) (code *Code) {
	if result := config.Database.First(&code, "user_id = ? AND type = ?", u.ID, "email"); result.Error != nil {
		code_gen, _ := GenerateCode(6)

		code = &Code{
			UserID:       u.ID,
			Type:         code_type,
			Code:         code_gen,
			AttemptCount: 0,
			ValidatedAt:  sql.NullTime{},
			ExpiredAt:    time.Now().Add(30 * time.Minute),
		}

		config.Database.Create(&code)

		return
	}

	if reset {
		code.Reset()
	}

	return
}

func (u *User) UpdatePassword(password string) {
	u.Password = string(services.EncryptPassword([]byte(password)))

	config.Database.Save(&u)
}

func (u *User) DecryptedPassword() string {
	return string(services.DecryptPassword([]byte(u.Password)))
}

func (u *User) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		UID:       u.UID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Bio: null.String{
			String: u.Bio.String,
			Valid:  u.Bio.Valid,
		},
		State:     u.State,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
