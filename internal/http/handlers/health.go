package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "application is fine",
		})
	}
}
