package services

import (
	"cambio-brl-usd/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type apiResponse struct {
	Base    string             `json:"base"`
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
}

var (
	NewHTTPRequest = http.NewRequest
	HTTPClientDo   = (&http.Client{}).Do
	DynamoScan     = func(c *dynamodb.Client, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return c.Scan(context.TODO(), input)
	}
	UnmarshalList = attributevalue.UnmarshalListOfMaps
)

var PutItemFn = func(client *dynamodb.Client, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return client.PutItem(context.TODO(), input)
}

var GetSecretValueFn = func(svc *secretsmanager.Client, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return svc.GetSecretValue(context.TODO(), input)
}

var JSONUnmarshalFn = json.Unmarshal

var SecretsFetcher = BuscarAPIKeyDoFixer
var SaveCotacao = SalvarCotacaoNoDynamo

func BuscarUltimaCotacao() models.Cotacao {
	token := SecretsFetcher()

	if token == "" {
		return BuscarUltimaCotacaoMock()
	}

	url := "https://api.apilayer.com/fixer/latest?base=BRL&symbols=USD"

	req, err := NewHTTPRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return BuscarUltimaCotacaoMock()
	}

	req.Header.Add("apikey", token)

	resp, err := HTTPClientDo(req)
	if err != nil {
		fmt.Println("Erro ao buscar cotação:", err)
		return BuscarUltimaCotacaoMock()
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return BuscarUltimaCotacaoMock()
	}

	if !apiResp.Success {
		fmt.Println("API retornou sucesso=false")
		return BuscarUltimaCotacaoMock()
	}

	usdRate := apiResp.Rates["USD"]
	base := apiResp.Base

	cotacao := models.Cotacao{
		MoedaOrigem:  base,
		MoedaDestino: "USD",
		Valor:        usdRate,
		DataHora:     time.Now(),
	}

	SaveCotacao(cotacao)

	return cotacao

}

func BuscarHistorico(inicio, fim time.Time) []models.Cotacao {

	cfg := carregarConfigAWS()

	client := dynamodb.NewFromConfig(cfg)

	// Convertendo datas para strings ISO
	dataInicio := inicio.Format(time.RFC3339)
	dataFim := fim.Format(time.RFC3339)

	// Filtro em data_hora
	filtro := expression.Name("data_hora").Between(expression.Value(dataInicio), expression.Value(dataFim))

	expr, err := expression.NewBuilder().WithFilter(filtro).Build()
	if err != nil {
		fmt.Println("Erro ao construir expressão:", err)
		return nil
	}

	// Scan com filtro - Tipo JPA Specifications
	input := &dynamodb.ScanInput{
		TableName:                 aws.String("Cotacoes"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := DynamoScan(client, input)
	if err != nil {
		fmt.Println("Erro ao fazer scan no DynamoDB:", err)
		return nil
	}

	var cotacoes []models.Cotacao
	if err := UnmarshalList(result.Items, &cotacoes); err != nil {
		fmt.Println("Erro ao converter resultados:", err)
		return nil
	}

	return cotacoes
}

func SalvarCotacaoNoDynamo(cotacao models.Cotacao) {

	cfg := carregarConfigAWS()

	client := dynamodb.NewFromConfig(cfg)

	item, err := attributevalue.MarshalMap(cotacao)
	if err != nil {
		fmt.Println("Erro ao converter cotação para DynamoDB:", err)
		return
	}

	_, err = PutItemFn(client, &dynamodb.PutItemInput{
		TableName: aws.String("Cotacoes"),
		Item:      item,
	})

	if err != nil {
		fmt.Println("Erro ao salvar no DynamoDB:", err)
	} else {
		fmt.Println("Cotação salva no DynamoDB com sucesso!")
	}
}

func BuscarAPIKeyDoFixer() string {
	secretName := "fixer-api-key-dev"

	cfg := carregarConfigAWS()

	svc := secretsmanager.NewFromConfig(cfg)

	result, err := GetSecretValueFn(svc, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})

	if err != nil {
		fmt.Println("Erro ao obter segredo:", err)
		return ""
	}

	var parsed map[string]string
	if err := JSONUnmarshalFn([]byte(*result.SecretString), &parsed); err != nil {
		fmt.Println("Erro ao interpretar JSON do segredo:", err)
		return ""
	}

	return parsed["fixer_api_key"]
}

func carregarConfigAWS() aws.Config {
	region := "us-east-1"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Println("Erro ao carregar configuração da AWS:", err)
		os.Exit(1) // encerra a aplicação
	}
	return cfg
}

func BuscarUltimaCotacaoMock() models.Cotacao {
	return models.Cotacao{
		MoedaOrigem:  "BRL",
		MoedaDestino: "USD",
		Valor:        5.00,
		DataHora:     time.Now(),
	}

}
