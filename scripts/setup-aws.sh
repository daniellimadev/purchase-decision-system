#!/bin/bash

# Script to configure AWS resources (SQS and SNS)

set -e

REGION="${AWS_REGION:-us-east-1}"
QUEUE_NAME="purchase-events-queue"
TOPIC_NAME="purchase-decisions-topic"

echo "🚀 Configuring AWS resources..."

echo "Creating SQS queue: $QUEUE_NAME"
QUEUE_URL=$(aws sqs create-queue \
  --queue-name "$QUEUE_NAME" \
  --region "$REGION" \
  --attributes VisibilityTimeout=30,MessageRetentionPeriod=345600 \
  --query 'QueueUrl' \
  --output text)

echo "✅ SQS queue created: $QUEUE_URL"

echo "Creating SNS topic: $TOPIC_NAME"
TOPIC_ARN=$(aws sns create-topic \
  --name "$TOPIC_NAME" \
  --region "$REGION" \
  --query 'TopicArn' \
  --output text)

echo "✅ SNS topic created: $TOPIC_ARN"

read -p "Do you want to create an email subscription? (y/n): " CREATE_EMAIL
if [ "$CREATE_EMAIL" = "y" ]; then
  read -p "Enter the email address: " EMAIL
  aws sns subscribe \
    --topic-arn "$TOPIC_ARN" \
    --protocol email \
    --notification-endpoint "$EMAIL" \
    --region "$REGION"
  echo "✅ Email subscription created. Please check your email to confirm."
fi

echo ""
echo "📝 Atualize seu arquivo .env com:"
echo "SQS_QUEUE_URL=$QUEUE_URL"
echo "SNS_TOPIC_ARN=$TOPIC_ARN"
