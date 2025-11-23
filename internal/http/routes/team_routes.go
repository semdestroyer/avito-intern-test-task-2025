package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
)

func registerTeamRoutes(s *ServiceDependencies.ServiceDependencies) {
	teamRoutes := s.Router.Group("team/")
	teamRoutes.GET("/get", handlers.TeamGet())
	teamRoutes.POST("/add", handlers.TeamAdd())
}
