package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/receipt/:id/points", getReceipt)
	router.POST("/receipts/process", processReceipts)
	router.Run("localhost:8080")
}
