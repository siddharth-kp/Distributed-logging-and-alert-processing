package main

import (
	"dlpas/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	s := gin.New()
	logger.InitLogger("service2", logger.BackendRedis)
	Logger := logger.GetLogger()
	Logger.Info("Service 2 started successfully")
	s.GET("/", func(c *gin.Context) {
		Logger.Info("endpoint hit service 2")
		c.JSON(200, gin.H{
			"service": "service 2",
			"message": "endpoint hit",
		})
	})

	Logger.Info("Service 2 listening on :8083")
	s.Run(":8083")
}
