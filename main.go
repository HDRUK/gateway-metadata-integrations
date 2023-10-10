package main

import (
	"fmt"
	"hdruk/federated-metadata/pkg/push"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("can't read .env file. resorting to os variables\n")
	}

	// Run the Push Service in it's own thread
	go push.Run()

	// Spawn a Pull Service on a cron scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	// TODO - reinstate this once we have federations
	// to begin running.
	// scheduler.Every(10).Second().Do(func() {
	// 	pull.Run()
	// })

	scheduler.StartBlocking()
}
