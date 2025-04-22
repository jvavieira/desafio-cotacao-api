resource "aws_dynamodb_table" "cotacoes" {
  name         = "Cotacoes"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "data_hora"

  attribute {
    name = "data_hora"
    type = "S"
  }
}