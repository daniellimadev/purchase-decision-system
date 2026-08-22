package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/purchase-decision-system/internal/domain"
	"github.com/purchase-decision-system/internal/infrastructure/config"
	log "github.com/sirupsen/logrus"
)

type SQSClient struct {
	client   *sqs.SQS
	queueURL string
	config   *config.AWSConfig
}

func NewSQSClient(cfg *config.AWSConfig) (*SQSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &SQSClient{
		client:   sqs.New(sess),
		queueURL: cfg.SQSQueueURL,
		config:   cfg,
	}, nil
}

func (s *SQSClient) ReceiveMessages() ([]*domain.PurchaseEvent, []*sqs.Message, error) {
	result, err := s.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: aws.Int64(s.config.SQSMaxMessages),
		WaitTimeSeconds:     aws.Int64(s.config.SQSWaitTime),
		VisibilityTimeout:   aws.Int64(s.config.SQSVisibilityTimeout),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	if len(result.Messages) == 0 {
		return nil, nil, nil
	}

	events := make([]*domain.PurchaseEvent, 0, len(result.Messages))
	for _, msg := range result.Messages {
		var event domain.PurchaseEvent
		if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
			log.WithError(err).WithField("message_id", *msg.MessageId).Error("Failed to unmarshal message")
			continue
		}
		events = append(events, &event)
	}

	return events, result.Messages, nil
}

func (s *SQSClient) DeleteMessage(receiptHandle *string) error {
	_, err := s.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: receiptHandle,
	})

	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *SQSClient) DeleteMessages(messages []*sqs.Message) error {
	if len(messages) == 0 {
		return nil
	}

	entries := make([]*sqs.DeleteMessageBatchRequestEntry, 0, len(messages))
	for i, msg := range messages {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(fmt.Sprintf("msg-%d", i)),
			ReceiptHandle: msg.ReceiptHandle,
		})
	}

	_, err := s.client.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(s.queueURL),
		Entries:  entries,
	})

	if err != nil {
		return fmt.Errorf("failed to delete messages batch: %w", err)
	}

	return nil
}

func (s *SQSClient) GetQueueAttributes() (map[string]string, error) {
	result, err := s.client.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(s.queueURL),
		AttributeNames: aws.StringSlice([]string{
			"ApproximateNumberOfMessages",
			"ApproximateNumberOfMessagesNotVisible",
			"ApproximateNumberOfMessagesDelayed",
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get queue attributes: %w", err)
	}

	attributes := make(map[string]string)
	for k, v := range result.Attributes {
		attributes[k] = *v
	}

	return attributes, nil
}
