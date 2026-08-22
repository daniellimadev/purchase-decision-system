resource "aws_sns_topic" "purchase_decisions" {
  name = var.topic_name

  tags = merge(
    local.common_tags,
    {
      Name = var.topic_name
    }
  )
}

resource "aws_sns_topic_policy" "purchase_decisions" {
  arn = aws_sns_topic.purchase_decisions.arn

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "sqs.amazonaws.com"
        }
        Action   = "SNS:Publish"
        Resource = aws_sns_topic.purchase_decisions.arn
      }
    ]
  })
}
