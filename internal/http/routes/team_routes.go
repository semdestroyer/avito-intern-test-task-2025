package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func registerTeamRoutes(r *gin.RouterGroup) {
	teamRoutes := r.Group("team/")
	teamRoutes.GET("/get", handlers.TeamGet())
	teamRoutes.POST("/add", handlers.TeamAdd())
}
