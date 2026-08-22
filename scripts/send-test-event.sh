#!/bin/bash

# Script to send a test event to SQS

set -e

# Settings
QUEUE_URL="${SQS_QUEUE_URL:-https://sqs.us-east-1.amazonaws.com/123456789/purchase-events-queue}"
REGION="${AWS_REGION:-us-east-1}"

# Generate random UUID
PURCHASE_ID=$(uuidgen)

# Create test event
EVENT=$(cat <<EOF
{
  "purchase_id": "$PURCHASE_ID",
  "customer_id": "customer-$(date +%s)",
  "amount": $((RANDOM % 10000 + 100)).50,
  "currency": "BRL",
  "merchant": "Test Store",
  "payment_method": "credit_card",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "metadata": {
    "ip_address": "192.168.1.1",
    "device": "mobile",
    "test": true
  }
}
EOF
)

echo "Sending test event to SQS..."
echo "$EVENT" | jq .

aws sqs send-message \
  --queue-url "$QUEUE_URL" \
  --message-body "$EVENT" \
  --region "$REGION"

echo ""
echo "✅ Event successfully submitted!"
echo "Purchase ID: $PURCHASE_ID"
