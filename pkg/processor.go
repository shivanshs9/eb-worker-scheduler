package pkg

import (
	"time"

	"github.com/shivanshs9/eb-worker-scheduler/pkg/cron"
	"github.com/shivanshs9/eb-worker-scheduler/pkg/sqs"
)

func (app *AppCls) pushJobToQueue(jobId string, event cron.CronEvent) error {
	current := time.Now()
	options := sqs.SqsOptions{
		QueueUrl:   app.options.QueueUrl,
		RetryCount: app.options.RetryCount,

		AttrJobPath:          event.Api,
		AttrJobTaskName:      event.Name,
		AttrJobScheduledTime: current.String(),

		DedupId: jobId,
		Body:    "{}",
	}
	app.log.Debugf("[%v] Pushing to SQS with options: %v", jobId, options)
	msgId, err := app.sqsClient.PushMessage(options)
	if err != nil {
		return err
	}
	app.log.Infof("[%v] Pushed to SQS with Msg ID: \"%v\"", jobId, msgId)
	return nil
}
