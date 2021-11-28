package main

import (
	"os"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/workers/daemons"
)

func CreateWorker(id string) daemons.Worker {
	switch id {
	case "cron_job":
		return daemons.NewCronJobWorker()
	default:
		return nil
	}
}

func main() {
	config.InitializeConfig()

	ARVG := os.Args[1:]

	for _, id := range ARVG {
		config.Logger.Println("Start finex-daemon: " + id)

		worker := CreateWorker(id)

		worker.Run()
	}
}
