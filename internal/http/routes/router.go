package router

import (
	"avito-intern-test-task-2025/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterPullRequestRoutes(r *gin.RouterGroup, handler *handlers.PrHandler) {
	prRoutes := r.Group("/pullRequest")
	prRoutes.POST("/merge", handler.PullRequestsMerge())
	prRoutes.POST("/create", handler.PullRequestsCreate())
	prRoutes.POST("/reassign", handler.PullRequestsReassign())
}

func RegisterTeamRoutes(r *gin.RouterGroup, handler *handlers.TeamHandler) {
	teamRoutes := r.Group("/team")
	teamRoutes.GET("/get", handler.TeamGet())
	teamRoutes.POST("/add", handler.TeamAdd())
	teamRoutes.POST("/bulkDeactivate", handler.TeamBulkDeactivate())
}

func RegisterUserRoutes(r *gin.RouterGroup, handler *handlers.UserHandler) {
	userRoutes := r.Group("/users")
	userRoutes.POST("/setIsActive", handler.UserSetIsActive())
	userRoutes.GET("/getReview", handler.UserGetReview())
}

func RegisterStatsRoutes(r *gin.RouterGroup, handler *handlers.StatsHandler) {
	statsRoutes := r.Group("/stats")
	statsRoutes.GET("/assignments", handler.GetAssignmentsStats())
}
