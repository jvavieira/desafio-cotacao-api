variable "region" {
  description = "AWS region"
  default     = "us-east-1"
}

variable "fixer_api_key_value" {
  description = "Chave da API do Fixer.io"
  type        = string
  sensitive   = true
}