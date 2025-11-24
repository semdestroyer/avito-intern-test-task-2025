package router

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPullRequestRoutes(r *gin.Engine, handler *handlers.PrHandler) {
	teamRoutes := r.Group("pullrequests/")
	teamRoutes.POST("/merge", handlers.PullRequestsMerge())
	teamRoutes.POST("/create", handlers.PullRequestsCreate())
	teamRoutes.POST("/reassign", handlers.PullRequestsReassign())
}

func RegisterTeamRoutes(r *gin.Engine, handler *handlers.TeamHandler) {
	teamRoutes := r.Group("team/")
	teamRoutes.GET("/get", handlers.TeamGet())
	teamRoutes.POST("/add", handlers.TeamAdd())
}

func RegisterUserRoutes(r *gin.Engine, handler *handlers.UserHandler) {
	userRoutes := r.Group("users/")
	userRoutes.POST("setIsActive", handler.UserSetIsActive())
	userRoutes.GET("getReview", handlers.UserGetReview())
}
