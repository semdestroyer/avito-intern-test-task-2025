package routes

import (
	"avito-intern-test-task-2025/internal/app"
	"avito-intern-test-task-2025/internal/http/handlers"
)

func registerTeamRoutes(s *app.Server) {
	teamRoutes := s.Router.Group("team/")
	teamRoutes.GET("/get", handlers.TeamGet())
	teamRoutes.POST("/add", handlers.TeamAdd())
}
