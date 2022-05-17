package cron

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

type JobProcessor interface {
	processJob(CronEvent) error
}

type CronEvent struct {
	Name    string `yaml:"name"`
	Crontab string `yaml:"schedule"`
	Api     string `yaml:"url"`
}

type Scheduler struct {
	gocron.Scheduler
	log *logrus.Logger
}

func NewScheduler(log *logrus.Logger) *Scheduler {
	return &Scheduler{
		*gocron.NewScheduler(time.UTC),
		log,
	}
}

func (scheduler *Scheduler) ScheduleCrons(crons []CronEvent, processor JobProcessor) error {
	for _, cron := range crons {
		if _, err := scheduler.Cron(cron.Crontab).Do(processor.processJob, cron); err != nil {
			return errors.New(fmt.Sprintf("Failed to schedule %v event: %v", cron, err))
		}
	}
	return nil
}
