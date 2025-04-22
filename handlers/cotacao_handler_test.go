package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUltimaCotacaoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	UltimaCotacao(c)

	// Verifica se o status de resposta foi 200 OK
	assert.Equal(t, 200, w.Code)
}

func TestHistoricoCotacaoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Caso de sucesso: datas válidas
	c.Request, _ = http.NewRequest("GET", "/cotacao/historico?inicio=2025-04-01T00:00&fim=2025-04-30T23:59", nil)
	HistoricoCotacao(c)
	assert.Equal(t, 200, w.Code)

}
