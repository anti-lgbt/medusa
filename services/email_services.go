package services

import (
	"encoding/json"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/types"
)

func SendEmail(key, language string, record map[string]interface{}) error {
	payload := &types.EngineMailerPayload{
		Key:      key,
		Language: language,
		Record:   record,
	}

	buf, _ := json.Marshal(payload)

	return config.Nats.Publish("engines:mailer", buf)
}
