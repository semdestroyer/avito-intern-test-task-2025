package routes

import (
	"avito-intern-test-task-2025/internal/app"
	"avito-intern-test-task-2025/internal/http/handlers"
)

func registerPullRequestRoutes(s *app.Server) {
	teamRoutes := s.Router.Group("pullrequests/")
	teamRoutes.POST("/merge", handlers.PullRequestsMerge())
	teamRoutes.POST("/create", handlers.PullRequestsCreate())
	teamRoutes.POST("/reassign", handlers.PullRequestsReassign())
}
