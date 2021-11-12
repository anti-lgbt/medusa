package models

import (
	"crypto/rand"
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/services"
)

type Code struct {
	ID           int64        `gorm:"primaryKey"`
	UserID       int64        `gorm:"type:bigint;not null;uniqueIndex:idx_user_id_and_type"`
	Type         string       `gorm:"type:character varying(10);not null;uniqueIndex:idx_user_id_and_type"`
	Code         string       `gorm:"type:character varying(6);not null"`
	AttemptCount int64        `gorm:"type:integer;not null"`
	ValidatedAt  sql.NullTime `gorm:"type:timestamp(0)"`
	ExpiredAt    time.Time    `gorm:"type:timestamp(0);not null"`
	CreatedAt    time.Time    `gorm:"type:timestamp(0);not null"`
	UpdatedAt    time.Time    `gorm:"type:timestamp(0);not null"`
	User         *User
}

func GenerateCode(length int) (string, error) {
	codeChars := "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(codeChars)
	for i := 0; i < length; i++ {
		buffer[i] = codeChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func (c *Code) OutAttempt() bool {
	max_attempt, _ := strconv.ParseInt(os.Getenv("MAX_ATTEMPT_COUNT"), 10, 64)

	return c.AttemptCount >= max_attempt
}

func (c *Code) Validated() bool {
	return c.ValidatedAt.Valid
}

func (c *Code) Expired() bool {
	return time.Now().After(c.ExpiredAt)
}

func (c *Code) Reset() error {
	code, _ := GenerateCode(6)
	c.Code = code
	c.ValidatedAt = sql.NullTime{}
	c.ExpiredAt = time.Now().Add(30 * time.Minute)

	config.Database.Save(&c)

	return nil
}

func (c *Code) Validation() {
	c.ValidatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	config.Database.Save(&c)
}

func (c *Code) SendCode(template string, language string) {
	services.SendEmail(
		template,
		language, map[string]interface{}{
			"code": c.Code,
		},
	)
}
