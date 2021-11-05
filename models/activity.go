package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
)

type Activity struct {
	ID            int64          `gorm:"primaryKey"`
	UserID        int64          `gorm:"bigint;not null;index"`
	Category      string         `gorm:"character varying(20);not null"`
	UserIP        string         `gorm:"character varying(50);not null"`
	UserIPCountry string         `gorm:"character varying(50)"`
	UserAgent     string         `gorm:"character varying;not null"`
	Topic         string         `gorm:"character varying(10);not null"`
	Action        string         `gorm:"character varying(10);not null"`
	Result        string         `gorm:"character varying(10);not null"`
	Data          sql.NullString `gorm:"json"`
	CreatedAt     time.Time      `gorm:"type:timestamp(0);not null;index"`
	UpdatedAt     time.Time      `gorm:"type:timestamp(0);not null;index"`
	User          *User
}

func CreateActivity(user_id int64, category string, user_ip, user_ip_country, user_agent, topic, action, result string, data sql.NullString) (*Activity, error) {
	activity := &Activity{
		UserID:        user_id,
		Category:      category,
		UserIP:        user_ip,
		UserIPCountry: user_ip_country,
		UserAgent:     user_agent,
		Topic:         topic,
		Action:        action,
		Result:        result,
		Data:          data,
	}

	r := config.Database.Create(&activity)

	return activity, r.Error
}
