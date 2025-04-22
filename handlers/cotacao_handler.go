package handlers

import (
    "net/http"
	"time"
    "github.com/gin-gonic/gin"
    "cambio-brl-usd/services"
)

func UltimaCotacao(c *gin.Context) {
    cotacao := services.BuscarUltimaCotacao()
    c.JSON(http.StatusOK, cotacao)
}

func HistoricoCotacao(c *gin.Context) {
    inicioStr := c.Query("inicio")
    fimStr := c.Query("fim")

    layout := "2006-01-02T15:04" // formato de entrada: "2025-04-18T18:30"

    inicio, err := time.Parse(layout, inicioStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"erro": "Data de início inválida"})
        return
    }

    fim, err := time.Parse(layout, fimStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"erro": "Data de fim inválida"})
        return
    }

    historico := services.BuscarHistorico(inicio, fim)
    c.JSON(http.StatusOK, historico)
}