package pkg

import (
	"time"

	"github.com/shivanshs9/eb-worker-scheduler/pkg/cron"
	"github.com/shivanshs9/eb-worker-scheduler/pkg/sqs"
)

func (app *AppCls) pushJobToQueue(event cron.CronEvent) error {
	current := time.Now()
	options := sqs.SqsOptions{
		QueueUrl:               app.options.QueueUrl,
		DeduplicationBufferSec: app.options.DeduplicationBufferSec,
		RetryCount:             app.options.RetryCount,

		AttrJobPath:          event.Api,
		AttrJobTaskName:      event.Name,
		AttrJobScheduledTime: current.String(),
	}
	msgId, err := app.sqsClient.PushMessage(options)
	if err != nil {
		return err
	}
	app.log.Infof("[%v] Pushed to SQS with Msg ID: %v", event.Name, msgId)
	return nil
}
