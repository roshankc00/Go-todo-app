package main

import (
	"os"

	routes "github.com/roshankc00/Go-todo-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port:= os.Getenv("PORT")
	if port ==""{
		port="8000"
	}
	router:=gin.New()
	router.Use(gin.Logger())
    routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.GET("/api",func(c *gin.Context){
		c.JSON(200,gin.H{"sucess":"Api is up and running "})
	})
	router.GET("/api-2",func(c *gin.Context){
		c.JSON(200,gin.H{"sucess":"Api is up and running "})
	})

	router.Run(":"+ port)
	
}