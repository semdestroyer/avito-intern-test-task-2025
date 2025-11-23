package routes

import (
	"avito-intern-test-task-2025/internal/app"
	"avito-intern-test-task-2025/internal/http/handlers"
)

func registerUserRoutes(s *app.Server) {
	userRoutes := s.Router.Group("users/")
	userRoutes.POST("setIsActive", handlers.UserSetIsActive())
	userRoutes.GET("getReview", handlers.UserGetReview())
}
