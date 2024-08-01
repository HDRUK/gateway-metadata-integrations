package main

import (
	"hdruk/federated-metadata/pkg/pull"
	"hdruk/federated-metadata/pkg/push"
	"hdruk/federated-metadata/pkg/utils"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.WriteGatewayAudit("can't read .env file. resorting to OS variables", "CONFIG")
	}

	debugLogs, err := strconv.Atoi(os.Getenv("DEBUG_LOGGING"))
	if err != nil {
		debugLogs = 0 // could not load config, err on side of fewer logs
	}

	if (debugLogs == 1) {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// Run the Push Service in it's own thread
	go push.Run()

	// Spawn a Pull Service on a cron scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	// TODO - reinstate this once we have federations
	// to begin running.
	scheduler.Every(1).Minute().Do(func() {
		pull.Run()
	})

	scheduler.StartBlocking()
}
