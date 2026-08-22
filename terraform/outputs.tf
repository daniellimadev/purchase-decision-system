output "sqs_queue_url" {
  description = "SQS queue URL"
  value       = aws_sqs_queue.purchase_events.url
}

output "sqs_queue_arn" {
  description = "SQS queue ARN"
  value       = aws_sqs_queue.purchase_events.arn
}

output "sqs_dlq_url" {
  description = "DLQ URL"
  value       = aws_sqs_queue.purchase_events_dlq.url
}

output "sns_topic_arn" {
  description = "SNS Topic ARN"
  value       = aws_sns_topic.purchase_decisions.arn
}
