package routes

import (
	"avito-intern-test-task-2025/internal/app"
	"avito-intern-test-task-2025/internal/http/handlers"
)

func RegisterRoutes(s *app.Server) {
	v1 := s.Router.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())
	registerPullRequestRoutes(s)
	registerTeamRoutes(s)
	registerUserRoutes(s)
}
