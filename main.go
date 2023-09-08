package main

import (
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/push"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't read .env file. aborting")
	}

	// Run the Push Service in it's own thread
	go push.Run()

	// Spawn a Pull Service on a cron scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	scheduler.Every(10).Second().Do(func() {
		pull.Run()
	})

	scheduler.StartBlocking()
}
