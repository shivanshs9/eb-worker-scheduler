package pkg

import (
	"fmt"
	"time"

	"github.com/shivanshs9/eb-worker-scheduler/pkg/cron"
	"github.com/shivanshs9/eb-worker-scheduler/pkg/sqs"
	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
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
	start := time.Now()
	// to track execution time
	defer func(start time.Time) {
		app.log.Infof("[%v] Took %v", event.Name, time.Since(start))
	}(start)

	idTimeU := fmt.Sprint(start.Round(60 * time.Second).Unix())
	idNameU := event.Name

	uniqueId := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(idTimeU+idNameU))

	app.log.Infof("[%v] Trigerred execution for \"%v\" at %v", uniqueId, event.Name, time.Now())
	if err := app.pushJobToQueue(uniqueId.String(), event); err != nil {
		app.log.Errorf("[%v] Failed to push to queue: %v", uniqueId, err)
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
	app.log.Infof("Starting the scheduler with %d crons", len(app.scheduler.Jobs()))
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
