package task

import (
	"github.com/robfig/cron/v3"
	"log"
)

func Start() {
	c := cron.New(cron.WithSeconds())
	spec := "*/5 * * * * *" //
	// task 1
	c.AddFunc(spec, func() {
		log.Printf("task 1")
	})
	// task 2
	c.AddFunc("*/1 * * * * *", func() {
		log.Printf("task 2")
	})
	c.Start()
}
