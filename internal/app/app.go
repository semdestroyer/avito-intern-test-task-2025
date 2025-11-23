package app

import (
	"avito-intern-test-task-2025/internal/config"
	"avito-intern-test-task-2025/internal/http/routes"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
	"avito-intern-test-task-2025/pkg/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {

	c := config.LoadConfig()
	db, err := db.InitDB(c)
	r := gin.Default()
	s := &ServiceDependencies.ServiceDependencies{
		Router: r,
		DB:     db,
	}

	routes.RegisterRoutes(s)

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application started")

}
