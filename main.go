package main

import (
	"hdruk/federated-metadata/pkg/push"
	"hdruk/federated-metadata/pkg/utils"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.WriteGatewayAudit("can't read .env file. resorting to OS variables", "CONFIG")
	}

	// Run the Push Service in it's own thread
	go push.Run()

	// Spawn a Pull Service on a cron scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	// TODO - reinstate this once we have federations
	// to begin running.
	// scheduler.Every(1).Minute().Do(func() {
	// 	pull.Run()
	// })

	scheduler.StartBlocking()
}
