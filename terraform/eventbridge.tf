resource "aws_cloudwatch_event_rule" "cotacao_agendada" {
  name                = "cotacao-agendada"
  schedule_expression = "cron(0 8,14,20 * * ? *)"
}

resource "aws_cloudwatch_event_target" "cotacao_lambda_target" {
  rule      = aws_cloudwatch_event_rule.cotacao_agendada.name
  target_id = "cotacao-lambda"
  arn       = aws_lambda_function.cotacao_lambda.arn
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cotacao_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cotacao_agendada.arn
}
