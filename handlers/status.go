package handlers

import (
	"github.com/gin-gonic/gin"

	"grpc-test/client"
)

func GetStatus(c *gin.Context) {
	grpcClient := client.NewClient(c)
	msg, ok := grpcClient.GetStatus()
	if ok {
		c.JSON(200, gin.H{"msg": msg})
	} else {
		c.JSON(500, gin.H{"msg": "sum ting fuked"})
	}
}
