package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/roshankc00/Go-todo-app/controllers"
	"github.com/roshankc00/Go-todo-app/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
  incomingRoutes.Use(middleware.Authenticate())
  incomingRoutes.GET("/users",controller.GetUsers())
  incomingRoutes.GET("/users/:user_id",controller.GetUser())
}
