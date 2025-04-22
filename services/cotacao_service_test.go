package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuscarUltimaCotacao(t *testing.T) {
	cotacao := BuscarUltimaCotacao()
	assert.Equal(t, "BRL", cotacao.MoedaOrigem)
	assert.Equal(t, "USD", cotacao.MoedaDestino)
	assert.Greater(t, cotacao.Valor, 0.0)
}
