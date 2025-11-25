package app

import (
	"avito-intern-test-task-2025/internal/config"
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/handlers"
	router "avito-intern-test-task-2025/internal/http/routes"
	"avito-intern-test-task-2025/internal/usecase"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
	"avito-intern-test-task-2025/pkg/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Run() {

	c := config.LoadConfig()
	database, err := db.InitDB(c)
	database.RunMigrations()
	r := gin.Default()
	s := &ServiceDependencies.ServiceDependencies{
		DB: database,
	}

	v1 := r.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())

	tr := repo.NewTeamRepo(s.DB)
	ur := repo.NewUserRepo(s.DB)
	pr := repo.NewPrRepo(s.DB, ur)
	uc := usecase.NewUserUsecase(ur, pr)
	tc := usecase.NewTeamUsecase(tr, ur)
	pc := usecase.NewPullRequestUsecase(ur, pr)
	uh := handlers.NewUserHandler(uc)
	th := handlers.NewTeamHandler(tc)
	ph := handlers.NewPrHandler(pc)

	router.RegisterUserRoutes(api, uh)
	router.RegisterTeamRoutes(api, th)
	router.RegisterPullRequestRoutes(api, ph)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application started")

}
