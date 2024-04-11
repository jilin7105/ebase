package task

import (
	"github.com/go-co-op/gocron"
	"time"
)

func InitTaskServer() *gocron.Scheduler {

	s := gocron.NewScheduler(time.Local)

	return s
}
