resource "aws_sqs_queue" "purchase_events" {
  name                       = var.queue_name
  visibility_timeout_seconds = var.visibility_timeout
  message_retention_seconds  = var.message_retention_period
  max_message_size          = var.max_message_size
  delay_seconds             = 0
  receive_wait_time_seconds = 20

  tags = merge(
    local.common_tags,
    {
      Name = var.queue_name
    }
  )
}

resource "aws_sqs_queue" "purchase_events_dlq" {
  name = "${var.queue_name}-dlq"

  tags = merge(
    local.common_tags,
    {
      Name = "${var.queue_name}-dlq"
    }
  )
}

resource "aws_sqs_queue_redrive_policy" "purchase_events" {
  queue_url = aws_sqs_queue.purchase_events.id

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.purchase_events_dlq.arn
    maxReceiveCount     = 3
  })
}
