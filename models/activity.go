package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
)

type Activity struct {
	ID            int64          `json:"id" gorm:"primaryKey"`
	UserID        int64          `json:"-"`
	Category      string         `json:"category"`
	UserIP        string         `json:"user_ip"`
	UserIPCountry string         `json:"user_ip_country"`
	UserAgent     string         `json:"user_agent"`
	Topic         string         `json:"topic"`
	Action        string         `json:"action"`
	Result        string         `json:"result"`
	Data          sql.NullString `json:"data"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
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
