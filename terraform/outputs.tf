output "ecr_repo_url" {
  value       = aws_ecr_repository.cotacao_api.repository_url
  description = "URL do repositório ECR"
}

output "ecr_repo_lambda_url" {
  description = "URL do repositório ECR da função Lambda"
  value       = aws_ecr_repository.cotacao_lambda.repository_url
}

output "dynamodb_table_name" {
  value       = aws_dynamodb_table.cotacoes.name
  description = "Nome da tabela DynamoDB"
}

output "secret_arn" {
  value       = aws_secretsmanager_secret.fixer_api_key.arn
  description = "ARN do segredo do Fixer API Key"
}

output "apprunner_service_url" {
  description = "URL pública do serviço App Runner"
  value       = aws_cloudformation_stack.apprunner_stack.outputs["ServiceUrl"]
}

output "lambda_function_name" {
  value = aws_lambda_function.cotacao_lambda.function_name
}

output "lambda_function_arn" {
  value = aws_lambda_function.cotacao_lambda.arn
}

output "eventbridge_rule" {
  value = aws_cloudwatch_event_rule.cotacao_agendada.schedule_expression
}