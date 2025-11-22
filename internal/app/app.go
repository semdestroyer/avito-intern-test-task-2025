package app

import (
	"avito-intern-test-task-2025/internal/http/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	r := gin.Default()
	routes.RegisterRoutes(r)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Application started")

}
