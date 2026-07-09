package main

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

const defaultCron = "0 * * * *"

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "characters" {
		if err := listCharacters(cfg); err != nil {
			log.Fatalf("characters: %v", err)
		}
		return
	}

	if len(os.Args) < 2 || os.Args[1] != "cron" {
		if err := syncStats(cfg); err != nil {
			log.Fatalf("sync: %v", err)
		}
		log.Println("stats updated")
		return
	}

	schedule := defaultCron
	if len(os.Args) > 2 && os.Args[2] != "" {
		schedule = os.Args[2]
	}

	c := cron.New()
	if _, err := c.AddFunc(schedule, func() {
		if err := syncStats(cfg); err != nil {
			log.Printf("sync: %v", err)
			return
		}
		log.Println("stats updated")
	}); err != nil {
		log.Fatalf("cron: %v", err)
	}

	log.Printf("running on schedule %q", schedule)
	c.Run()
}
