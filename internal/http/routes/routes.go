package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())
	registerPullRequestRoutes(api)
	registerTeamRoutes(api)
	registerUserRoutes(api)
}
