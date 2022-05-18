package sqs

import (
	"errors"
	"strings"
	"time"

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
	QueueUrl   string
	RetryCount int

	AttrJobPath          string
	AttrJobScheduledTime string
	AttrJobTaskName      string

	DedupId string
	Body    string
}

func NewSqsClient(log *logrus.Logger) *Client {
	mySession := session.Must(session.NewSession())
	return &Client{
		sqs: sqs.New(mySession),
		log: log,
	}
}

func (client *Client) PushMessage(options SqsOptions) (msgId string, err error) {
	dataType := "String"
	attributes := map[string]*sqs.MessageAttributeValue{
		"beanstalk.sqsd.path": {
			StringValue: &options.AttrJobPath,
			DataType:    &dataType,
		},
		"beanstalk.sqsd.task_name": {
			StringValue: &options.AttrJobTaskName,
			DataType:    &dataType,
		},
		"beanstalk.sqsd.scheduled_time": {
			StringValue: &options.AttrJobScheduledTime,
			DataType:    &dataType,
		},
	}
	input := &sqs.SendMessageInput{
		QueueUrl:          &options.QueueUrl,
		MessageAttributes: attributes,
		MessageBody:       &options.Body,
	}
	if isFifoQueue := strings.HasSuffix(options.QueueUrl, ".fifo"); isFifoQueue {
		// Config for FIFO queue
		input.MessageGroupId = &options.AttrJobPath
		input.MessageDeduplicationId = &options.DedupId
	}

	attempt := 1
	for {
		if attempt > options.RetryCount {
			err = errors.New("Max attempts reached")
			break
		}
		output, err := client.sqs.SendMessage(input)
		if err != nil {
			client.log.Warnf("[%v/%v] Failed to send message: %v", attempt, options.RetryCount, err)
			attempt += 1
			time.Sleep(2 * time.Second)
			continue
		}
		client.log.Debugf("Response: %v", output)
		msgId = *output.MessageId
		break
	}
	return
}
