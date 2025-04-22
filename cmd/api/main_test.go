package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// cria uma instância da rota para teste
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/cotacao/ultima", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensagem": "mock cotação"})
	})
	r.GET("/cotacao/historico", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensagem": "mock histórico"})
	})
	return r
}

func TestUltimaCotacaoEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cotacao/ultima", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "mock cotação")
}

func TestHistoricoCotacaoEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cotacao/historico", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "mock histórico")
}
