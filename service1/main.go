package main

import (
	"dlpas/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	s := gin.New()
	logger.InitLogger("service1")
	Logger := logger.GetLogger()
	Logger.Info("Service 1 started successfully")
	s.GET("/", func(c *gin.Context) {
		Logger.Info("endpoint hit service 1")
		c.JSON(200, gin.H{
			"service": "service 1",
			"message": "endpoint hit",
		})
	})

	Logger.Info("Service 1 listening on :8081")
	s.Run(":8081")
}
