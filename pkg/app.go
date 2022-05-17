package pkg

import (
	"github.com/shivanshs9/eb-worker-scheduler/pkg/cron"
	"github.com/shivanshs9/eb-worker-scheduler/pkg/sqs"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	sqs.ReceiveMessageOptions
	YamlPath string
}

type AppCls struct {
	sqsClient *sqs.Client
	scheduler *cron.Scheduler
	options   *AppOptions
	log       *logrus.Logger
}

// func (app *AppCls) processJob(job *sqs.SQSJobRequest) error {
// 	// to track execution time
// 	defer func(start time.Time) {
// 		app.log.Infof("[%v] Took %v", job.SqsMsgId, time.Since(start))
// 	}(time.Now())

// 	app.log.Infof("[%v] Sending POST to %v", job.SqsMsgId, job.AttrJobPath)
// 	resp, err := app.httpClient.PostRequest(*job)
// 	if resp != nil {
// 		defer resp.Body.Close()
// 	}
// 	if err != nil {
// 		return err
// 	} else if resp.StatusCode != 200 {
// 		return errors.New(fmt.Sprintf("Received %v from the API call", resp.Status))
// 	}
// 	return nil
// }

func (app *AppCls) start() {
	app.log.Infof("Starting the app with following config: %v", *app.options)
	cronData, err := cron.ParseYaml(app.options.YamlPath)
	if err != nil {
		app.log.Fatalf("Failed parsing cron file %v: %v", app.options.YamlPath, err)
	}
	if err = app.scheduler.ScheduleCrons(cronData.Crons); err != nil {
		app.log.Fatalf("Failed to schedule crons: %v", err)
	}
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
