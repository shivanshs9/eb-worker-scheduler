package sqs

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

const MAX_ERROR_IGNORE = 5

type Client struct {
	sqs *sqs.SQS
	log *logrus.Logger
}

type SqsOptions struct {
	QueueUrl               string
	DeduplicationBufferSec int
	RetryCount             int

	AttrJobPath          string
	AttrJobScheduledTime string
	AttrJobTaskName      string
}

func NewSqsClient(log *logrus.Logger) *Client {
	mySession := session.Must(session.NewSession())
	return &Client{
		sqs: sqs.New(mySession),
		log: log,
	}
}

func (client *Client) PushMessage(options SqsOptions) (msgId string, err error) {
	// dedupId := time.Now().Format("")
	attributes := map[string]*sqs.MessageAttributeValue{
		"beanstalk.sqsd.path": {
			StringValue: &options.AttrJobPath,
		},
		"beanstalk.sqsd.task_name": {
			StringValue: &options.AttrJobTaskName,
		},
		"beanstalk.sqsd.scheduled_time": {
			StringValue: &options.AttrJobScheduledTime,
		},
	}
	body := "{}"
	input := &sqs.SendMessageInput{
		QueueUrl: &options.QueueUrl,
		// MessageDeduplicationId: &dedupId,
		MessageAttributes: attributes,
		MessageBody:       &body,
	}
	output, err := client.sqs.SendMessage(input)
	if err != nil {
		return
	}
	msgId = *output.MessageId
	return
}
