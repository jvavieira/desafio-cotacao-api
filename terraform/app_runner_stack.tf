resource "aws_cloudformation_stack" "apprunner_stack" {
  name = "cotacao-api-runner"

template_body = <<TEMPLATE
AWSTemplateFormatVersion: '2010-09-09'
Resources:
  CotacaoApiService:
    Type: AWS::AppRunner::Service
    Properties:
      ServiceName: cotacao-api-service
      SourceConfiguration:
        ImageRepository:
          ImageIdentifier: 087791688156.dkr.ecr.us-east-1.amazonaws.com/cotacao-api:latest
          ImageRepositoryType: ECR
          ImageConfiguration:
            Port: "8080"
        AutoDeploymentsEnabled: true
        AuthenticationConfiguration:
          AccessRoleArn: ${aws_iam_role.apprunner_ecr_access.arn}
      InstanceConfiguration:
        Cpu: "1 vCPU"
        Memory: "2 GB"
        InstanceRoleArn: ${aws_iam_role.apprunner_exec_role.arn}

Outputs:
  ServiceUrl:
    Description: "App Runner Public URL"
    Value: !GetAtt CotacaoApiService.ServiceUrl
TEMPLATE
}