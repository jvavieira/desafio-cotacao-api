variable "fixer_api_key_value" {
  description = "Chave da API do Fixer.io"
  type        = string
  sensitive   = true
}

resource "aws_secretsmanager_secret" "fixer_api_key" {
  name = "fixer-api-key-dev"
}

resource "aws_secretsmanager_secret_version" "fixer_api_key_value" {
  secret_id     = aws_secretsmanager_secret.fixer_api_key.id
  secret_string = jsonencode({
    fixer_api_key = var.fixer_api_key_value
  })
}
