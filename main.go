package main

import (
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/roshankc00/Go-todo-app/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.TodoRoutes(router)

	router.Run(":" + port)
}
