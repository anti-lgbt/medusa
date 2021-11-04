package models

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models/datatypes"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
)

type UserData struct {
	Language string   `json:"language"`
	IPS      []string `json:"ips"`
}

type User struct {
	ID        int64                `json:"id" gorm:"primaryKey"`
	UID       string               `json:"uid" gorm:"type:character varying(20);not null"`
	Email     string               `json:"email" gorm:"type:character varying(50);not null;uniqueIndex"`
	Password  string               `json:"-" gorm:"type:text;not null"`
	FirstName string               `json:"first_name" gorm:"type:character varying(50);not null;index"`
	LastName  string               `json:"last_name" gorm:"type:character varying(50);not null;index"`
	Bio       datatypes.NullString `json:"bio" gorm:"type:text"`
	State     types.UserState      `json:"state" gorm:"type:character varying(10):not null;index"`
	Role      types.UserRole       `json:"role" gorm:"type:character varying(10):not null;index"`
	Avatar    datatypes.NullString `json:"-" gorm:"type:text"`
	Data      sql.NullString       `json:"-" gorm:"type:text"`
	CreatedAt time.Time            `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
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
			ValidatedAt:  datatypes.NullTime{},
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
