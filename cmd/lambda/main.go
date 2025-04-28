package main

import (
	"cambio-brl-usd/services"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	services.BuscarUltimaCotacao()
	return "Ultima cotação buscada com sucesso!", nil
}

func main() {
	lambda.Start(handler)
}
