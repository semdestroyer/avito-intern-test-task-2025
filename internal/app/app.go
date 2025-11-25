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
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	pc := usecase.NewPullRequestUsecase(ur, pr)
	tc := usecase.NewTeamUsecase(tr, ur, pr, pc)
	sc := usecase.NewStatsUsecase(pr)
	uh := handlers.NewUserHandler(uc)
	th := handlers.NewTeamHandler(tc)
	ph := handlers.NewPrHandler(pc)
	sh := handlers.NewStatsHandler(sc)

	router.RegisterUserRoutes(api, uh)
	router.RegisterTeamRoutes(api, th)
	router.RegisterPullRequestRoutes(api, ph)
	router.RegisterStatsRoutes(api, sh)

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
