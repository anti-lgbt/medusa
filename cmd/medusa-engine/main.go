package main

import (
	"os"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/workers/engines"
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

		sub, _ := config.Nats.SubscribeSync("engines:" + id)

		for {
			m, err := sub.NextMsg(1 * time.Second)

			if err != nil {
				continue
			}

			// config.Logger.Infof("Receive message: %s", string(m.Data))
			if err := worker.Process(m.Data); err == nil {
				m.Ack()
			} else {
				config.Logger.Errorf("Worker error: %v", err.Error())
			}
		}
	}
}
