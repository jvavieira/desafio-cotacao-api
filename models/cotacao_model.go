package models

import "time"

type Cotacao struct {
    MoedaOrigem  string    `json:"moeda_origem" dynamodbav:"moeda_origem"`
    MoedaDestino string    `json:"moeda_destino" dynamodbav:"moeda_destino"`
    Valor        float64   `json:"valor" dynamodbav:"valor"`
    DataHora     time.Time `json:"data_hora" dynamodbav:"data_hora"`
}
