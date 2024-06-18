package utils

import (
	"github.com/robfig/cron"
	"log"
)

func init() {
	log.Default()
}

func CronJob() {
	c := cron.New()
	err := c.AddFunc("@hourly", UpdateDb)
	log.Println("CronJob started")
	if err != nil {
		return
	}
	c.Start()

}
