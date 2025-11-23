package routes

import (
	"avito-intern-test-task-2025/internal/http/handlers"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
)

func registerPullRequestRoutes(s *ServiceDependencies.ServiceDependencies) {
	teamRoutes := s.Router.Group("pullrequests/")
	teamRoutes.POST("/merge", handlers.PullRequestsMerge())
	teamRoutes.POST("/create", handlers.PullRequestsCreate())
	teamRoutes.POST("/reassign", handlers.PullRequestsReassign())
}
