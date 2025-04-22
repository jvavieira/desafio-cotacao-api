package main

import (
	"cambio-brl-usd/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/cotacao/ultima", handlers.UltimaCotacao)
	r.GET("/cotacao/historico", handlers.HistoricoCotacao)
	r.Run(":8080")
}
