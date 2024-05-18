package cron

import (
	"log/slog"
	"sync"
	"time"
)

// NewDailyJob constructs job which will  be triggered on the
// daily basis at the given point of time.
// NOTE: it expects time in UTC timezone. It's -3 hours compared to Kyiv Time:
// 12:00 UTC = 15:00 by Kyiv
func NewDailyJob(hour, min int, log *slog.Logger) *Job {
	now := time.Now()
	tickAt := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.UTC)
	if now.After(tickAt) {
		tickAt = tickAt.Add(time.Hour * 24)
	}

	tickAfter := tickAt.Sub(now)
	return &Job{
		interval: time.Hour * 24,
		ticker:   time.NewTicker(tickAfter),
		log:      log,
	}
}

func NewJob(interval time.Duration, log *slog.Logger) *Job {
	return &Job{
		interval: interval,
		ticker:   time.NewTicker(interval),
		log:      log,
	}
}

type Job struct {
	interval time.Duration
	ticker   *time.Ticker
	log      *slog.Logger
}

func (j *Job) Do(fn func() error) {
	var once sync.Once
	go func() {
		for range j.ticker.C {
			once.Do(func() {
				j.ticker.Reset(j.interval)
			})
			if err := fn(); err != nil {
				j.log.Error("Failed to do cron task", "err", err)
			}
		}
	}()
}
