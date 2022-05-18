package pkg

import (
	"time"

	"github.com/shivanshs9/eb-worker-scheduler/pkg/cron"
	"github.com/shivanshs9/eb-worker-scheduler/pkg/sqs"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	sqs.SqsOptions
	YamlPath string
}

type AppCls struct {
	sqsClient *sqs.Client
	scheduler *cron.Scheduler
	options   *AppOptions
	log       *logrus.Logger
}

var _ cron.JobProcessor = &AppCls{}

func (app *AppCls) ProcessJob(event cron.CronEvent) {
	// to track execution time
	defer func(start time.Time) {
		app.log.Infof("[%v] Took %v", event.Name, time.Since(start))
	}(time.Now())

	app.log.Infof("[%v] Trigerred execution", event.Name)
	if err := app.pushJobToQueue(event); err != nil {
		app.log.Errorf("[%v] Failed to push to queue: %v", event.Name, err)
	}
}

func (app *AppCls) start() {
	app.log.Infof("Starting the app with following config: %v", *app.options)
	cronData, err := cron.ParseYaml(app.options.YamlPath)
	if err != nil {
		app.log.Fatalf("Failed parsing cron file %v: %v", app.options.YamlPath, err)
	}
	if err = app.scheduler.ScheduleCrons(cronData.Crons, app); err != nil {
		app.log.Fatalf("Failed to schedule crons: %v", err)
	}
	app.log.Infof("Starting the scheduler")
	app.scheduler.StartBlocking()
}

func StartApp(options *AppOptions, log *logrus.Logger) {
	app := &AppCls{
		sqsClient: sqs.NewSqsClient(log),
		scheduler: cron.NewScheduler(log),
		log:       log,
		options:   options,
	}
	app.start()
}
