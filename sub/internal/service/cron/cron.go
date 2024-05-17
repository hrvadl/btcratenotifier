package cron

import (
	"fmt"
	"log/slog"
	"time"
)

const operation = "cron job"

func NewJob(interval time.Duration, log *slog.Logger) *Job {
	return &Job{
		ticker: *time.NewTicker(interval),
		log:    log,
	}
}

type Job struct {
	ticker time.Ticker
	log    *slog.Logger
}

func (j *Job) Do(fn func() error) chan error {
	errCh := make(chan error)
	go func() {
		for range j.ticker.C {
			j.log.Info("Doing the job")
			if err := fn(); err != nil {
				j.log.Error("failed to do cron task", "err", err)
				errCh <- fmt.Errorf("%s: run failed: %w", operation, err)
				return
			}
		}
	}()
	return errCh
}
