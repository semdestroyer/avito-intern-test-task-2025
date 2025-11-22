package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("users/")
	userRoutes.POST("setIsActive", handlers.UserSetIsActive())
	userRoutes.GET("getReview", handlers.UserGetReview())
}
