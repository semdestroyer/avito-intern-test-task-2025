package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
)

func RegisterRoutes(s *ServiceDependencies.ServiceDependencies) {
	v1 := s.Router.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())
	registerPullRequestRoutes(s)
	registerTeamRoutes(s)
	registerUserRoutes(s)
}
