package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/roshankc00/Go-todo-app/controllers"
	"github.com/roshankc00/Go-todo-app/middleware"
)

func TodoRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate()) 
	incomingRoutes.POST("/todos",controller.CreateTodo())
	incomingRoutes.GET("/todos",controller.GetTodos())
	incomingRoutes.GET("/todos/:todo_id",controller.GetTodo())
	incomingRoutes.PATCH("/todos/:todo_id",controller.UpdateTodo())
	incomingRoutes.DELETE("/todos/:todo_id",controller.DeleteTodo())
}
