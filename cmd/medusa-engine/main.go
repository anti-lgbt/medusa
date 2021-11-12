package main

import (
	"os"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/workers/engines"
	"github.com/nats-io/nats.go"
)

func CreateWorker(id string) engines.Worker {
	switch id {
	case "mailer":
		return engines.NewMailerEngineWorker()
	default:
		return nil
	}
}

func main() {
	config.InitializeConfig()

	ARVG := os.Args[1:]

	for _, id := range ARVG {
		config.Logger.Println("Start finex-engine: " + id)
		worker := CreateWorker(id)

		config.Nats.Subscribe(id, func(m *nats.Msg) {
			config.Logger.Infof("Receive message: %s", string(m.Data))

			if err := worker.Process(m.Data); err == nil {
				m.Ack()
			} else {
				config.Logger.Errorf("Worker error: %v", err.Error())
			}
		})
	}
}
