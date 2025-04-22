resource "aws_ecr_repository" "cotacao_api" {
  name                 = "cotacao-api"
  image_tag_mutability = "MUTABLE"

  encryption_configuration {
    encryption_type = "AES256"
  }
}

resource "aws_ecr_repository" "cotacao_lambda" {
  name                 = "cotacao-lambda"
  image_tag_mutability = "MUTABLE"

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name = "cotacao-lambda"
  }
}
