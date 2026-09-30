package cron

import (
	"context"
	"fmt"
	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/lib/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

// Job 由业务通过 fx.ResultTags(`group:"cron_jobs"`) 注册；模板不自带业务任务。
type Job struct {
	Name     string
	Schedule string
	Run      func(context.Context)
}
type Params struct {
	fx.In
	Jobs []Job `group:"cron_jobs"`
}

func NewCron(lifecycle fx.Lifecycle, cfg config.Config, params Params) error {
	if !cfg.Cron.On {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.PrintfLogger(log.Logger)), cron.Recover(cron.PrintfLogger(log.Logger))))
	for _, job := range params.Jobs {
		if job.Name == "" || job.Run == nil {
			cancel()
			return fmt.Errorf("cron job requires a name and callback")
		}
		if _, err := scheduler.AddFunc(job.Schedule, func() { job.Run(ctx) }); err != nil {
			cancel()
			return fmt.Errorf("register cron %s: %w", job.Name, err)
		}
	}
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error { scheduler.Start(); return nil }, OnStop: func(ctx context.Context) error {
		cancel()
		done := scheduler.Stop()
		select {
		case <-done.Done():
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}})
	return nil
}
