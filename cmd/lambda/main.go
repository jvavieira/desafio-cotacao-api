package main

import (
	"cambio-brl-usd/services"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) {
	services.BuscarUltimaCotacao()
}

func main() {
	lambda.Start(Handler)
}
