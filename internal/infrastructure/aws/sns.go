package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/purchase-decision-system/internal/domain"
	"github.com/purchase-decision-system/internal/infrastructure/config"
	log "github.com/sirupsen/logrus"
)

type SNSClient struct {
	client   *sns.SNS
	topicARN string
}

func NewSNSClient(cfg *config.AWSConfig) (*SNSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &SNSClient{
		client:   sns.New(sess),
		topicARN: cfg.SNSTopicARN,
	}, nil
}

func (s *SNSClient) PublishDecisionEvent(event *domain.DecisionEvent) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := string(messageBytes)

	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(s.topicARN),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("purchase_decision"),
			},
			"status": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event.Status)),
			},
			"decision_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event.DecisionType)),
			},
		},
	}

	result, err := s.client.Publish(input)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.WithFields(log.Fields{
		"message_id":  *result.MessageId,
		"purchase_id": event.PurchaseID,
		"status":      event.Status,
	}).Info("Decision event published to SNS")

	return nil
}

func (s *SNSClient) PublishBatchDecisionEvents(events []*domain.DecisionEvent) error {
	for _, event := range events {
		if err := s.PublishDecisionEvent(event); err != nil {
			log.WithError(err).WithField("purchase_id", event.PurchaseID).Error("Failed to publish event")
			continue
		}
	}
	return nil
}
