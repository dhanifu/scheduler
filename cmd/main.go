package main

import (
	"errors"
	"go-scheduler/config"
	"go-scheduler/logger"
	"go-scheduler/scheduler"
	"log"
	"os"
	"strings"
)

func main() {
	log.Printf("Initializing scheduler...")
	var env string
	if len(os.Args) > 1 {
		log.Printf("Environment provided: %v", os.Args[1])
		if strings.Contains(os.Args[1], "stag") {
			env = "staging"
		} else if strings.Contains(os.Args[1], "prod") {
			env = "production"
		} else {
			err := errors.New("Invalid environment provided: " + os.Args[1])
			log.Fatal(err)
		}
	} else {
		env = "local"
	}

	// init config
	conf := config.LoadConfig("./", env)
	_ = logger.InitZerolog(conf)

	// init cron
	logger.Info("Starting scheduler...")
	scheduler := scheduler.NewScheduler(conf)
	scheduler.Start()

	select {}
}
