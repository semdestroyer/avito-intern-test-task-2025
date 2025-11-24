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
)

func Run() {

	c := config.LoadConfig()
	database, err := db.InitDB(c)
	database.RunMigrations()
	r := gin.Default()
	s := &ServiceDependencies.ServiceDependencies{
		//TODO: подумать нужно ли оставить всю структуру в целом. это может быть полезно чтобы тащить из pkg но пока у нас лишь db
		DB: database,
	}

	v1 := r.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())

	ur := repo.NewUserRepo(s.DB)
	uc := usecase.NewUserUsecase(ur)
	uh := handlers.NewUserHandler(uc)
	router.RegisterUserRoutes(r, uh)

	//Остальные пути потом
	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application started")

}
