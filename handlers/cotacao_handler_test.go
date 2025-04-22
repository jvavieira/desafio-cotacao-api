package handlers_test

import (
	"cambio-brl-usd/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/cotacao/ultima", handlers.UltimaCotacao)
	r.GET("/cotacao/historico", handlers.HistoricoCotacao)
	return r
}

func TestUltimaCotacaoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handlers.UltimaCotacao(c)

	// Verifica se o status de resposta foi 200 OK
	assert.Equal(t, 200, w.Code)
}

func TestHistoricoCotacaoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Caso de sucesso: datas válidas
	c.Request, _ = http.NewRequest("GET", "/cotacao/historico?inicio=2025-04-01T00:00&fim=2025-04-30T23:59", nil)
	handlers.HistoricoCotacao(c)
	assert.Equal(t, 200, w.Code)

}

func TestUltimaCotacao(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cotacao/ultima", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "USD") // ou alguma string esperada no JSON
}

func TestHistoricoCotacao_ComParametrosValidos(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cotacao/historico?inicio=2025-01-01&fim=2025-01-10", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestHistoricoCotacao_DataInvalida(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cotacao/historico?inicio=invalid&fim=2025-01-10", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}
