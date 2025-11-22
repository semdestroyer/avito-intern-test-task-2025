package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func registerPullRequestRoutes(r *gin.RouterGroup) {
	teamRoutes := r.Group("pullrequests/")
	teamRoutes.POST("/merge", handlers.PullRequestsMerge())
	teamRoutes.POST("/create", handlers.PullRequestsCreate())
	teamRoutes.POST("/reassign", handlers.PullRequestsReassign())
}
