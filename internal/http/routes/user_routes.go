package routes

import (
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/handlers"
	"avito-intern-test-task-2025/internal/usecase"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
)

func registerUserRoutes(s *ServiceDependencies.ServiceDependencies) {
	//TODO: когда просплюсь понять мб нужно тащить ServiceDependencies до конца, а не вкладывать каждый сервис в сервис
	ur := repo.NewUserRepo(s.DB)
	uc := usecase.NewUserUsecase(ur)
	uh := handlers.NewUserHandler(uc)
	userRoutes := s.Router.Group("users/")
	userRoutes.POST("setIsActive", uh.UserSetIsActive())
	userRoutes.GET("getReview", handlers.UserGetReview())
}
