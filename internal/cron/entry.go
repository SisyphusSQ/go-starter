package cron

import (
	"time"

	"github.com/robfig/cron/v3"
	"go-starter/internal/lib/log"
)

func New() {
	log.Logger.Info("starting cron...")

	timezone, _ := time.LoadLocation("Asia/Shanghai")
	Cron := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	Cron.AddFunc("0 0 2 * * ?", cronSample)

	Cron.Start()
}
