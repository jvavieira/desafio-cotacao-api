resource "aws_lambda_function" "cotacao_lambda" {
  function_name = "cotacao-busca-salva"
  role          = aws_iam_role.lambda_exec_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.cotacao_lambda.repository_url}:latest"
  timeout       = 10
}
