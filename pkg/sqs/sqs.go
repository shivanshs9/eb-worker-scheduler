package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

const MAX_ERROR_IGNORE = 5

type Client struct {
	sqs *sqs.SQS
	log *logrus.Logger
}

type ReceiveMessageOptions struct {
	QueueUrl            string
	MaxBufferedMessages int
	RetryCount          int
	DefaultHttpPath     string
}

type SQSJobRequest struct {
	AttrJobPath          string
	AttrJobScheduledTime string
	AttrJobTaskName      string

	SqsMsgId           string
	SqsQueueUrl        string
	SqsFirstReceivedAt string

	Body string

	receiptHandle string
}

func (msg SQSJobRequest) String() string {
	return fmt.Sprintf("[%v]: %v [%v] [%v]", msg.SqsMsgId, msg.Body, msg.AttrJobTaskName, msg.AttrJobPath)
}

func NewSqsClient(log *logrus.Logger) *Client {
	mySession := session.Must(session.NewSession())
	return &Client{
		sqs: sqs.New(mySession),
		log: log,
	}
}

func (client *Client) ReceiveMessageStream(options ReceiveMessageOptions, stop chan struct{}) chan *SQSJobRequest {
	// errors will be logged and ignored
	client.log.Info("Starting the SQS Messages Stream")
	stream := make(chan *SQSJobRequest, options.MaxBufferedMessages)
	errorCnt := 0
	go func() {
		for {
			select {
			case <-stop: // triggered when the stop channel is closed
				client.log.Info("Stopping the Stream")
				return // exit
			default:
				msgs, err := client.receiveMessage(options)
				if err != nil {
					errorCnt++
					client.log.WithError(err).Warn("Received error from SQS Client")
					if errorCnt >= MAX_ERROR_IGNORE {
						client.log.Fatal("Max attempts reached, exiting...")
					}
				} else {
					errorCnt = 0
					client.log.WithField("NumMessages", len(msgs)).Info("Received messages")
					for _, msg := range msgs {
						stream <- msg
					}
				}
			}
		}
	}()
	return stream
}

func (client *Client) AcknowledgeMessage(job *SQSJobRequest) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &job.SqsQueueUrl,
		ReceiptHandle: &job.receiptHandle,
	}
	_, err := client.sqs.DeleteMessage(input)
	if err != nil {
		return err
	}
	client.log.Infof("[%v] Deleted message from the SQS", job.SqsMsgId)
	return nil
}

func (client *Client) receiveMessage(options ReceiveMessageOptions) (jobs []*SQSJobRequest, err error) {
	maxMsgCount := int64(options.MaxBufferedMessages)
	waitTime := int64(20)
	attributeName := "*"
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &options.QueueUrl,
		MaxNumberOfMessages:   &maxMsgCount,
		WaitTimeSeconds:       &waitTime,
		MessageAttributeNames: []*string{&attributeName},
	}
	output, err := client.sqs.ReceiveMessage(input)
	if err != nil {
		return
	}
	for _, msg := range output.Messages {
		client.log.Debugf("Received message, %v", msg)
		jobPath := options.DefaultHttpPath
		if path, found := msg.MessageAttributes["beanstalk.sqsd.path"]; found {
			jobPath = *path.StringValue
		}
		attrScheduledTime := ""
		if time, found := msg.MessageAttributes["beanstalk.sqsd.scheduled_time"]; found {
			attrScheduledTime = *time.StringValue
		}
		attrTaskName := ""
		if taskName, found := msg.MessageAttributes["beanstalk.sqsd.task_name"]; found {
			attrTaskName = *taskName.StringValue
		}
		job := &SQSJobRequest{
			SqsMsgId:    *msg.MessageId,
			SqsQueueUrl: options.QueueUrl,

			AttrJobPath:          jobPath,
			AttrJobScheduledTime: attrScheduledTime,
			AttrJobTaskName:      attrTaskName,

			Body:          *msg.Body,
			receiptHandle: *msg.ReceiptHandle,
		}
		jobs = append(jobs, job)
	}
	return
}
