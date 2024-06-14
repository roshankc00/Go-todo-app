package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roshankc00/Go-todo-app/database"
	"github.com/roshankc00/Go-todo-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var todoCollection *mongo.Collection=database.OpenCollection(database.Client,"todo")



func GetTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        todoID := c.Param("todo_id")

        objID, err := primitive.ObjectIDFromHex(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo_id format"})
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var todo models.Todo
        err = todoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "data":    todo,
        })
    }
}


func GetTodos() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        recordPerPage, err := strconv.Atoi(c.DefaultQuery("recordPerPage", "10"))
        if err != nil || recordPerPage < 1 {
            recordPerPage = 10
        }

        page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
        if err != nil || page < 1 {
            page = 1
        }

        startIndex := (page - 1) * recordPerPage

        matchStage := bson.D{{"$match", bson.D{{}}}}
        groupStage := bson.D{{"$group", bson.D{
            {"_id", bson.D{{"_id", "null"}}},
            {"total_count", bson.D{{"$sum", 1}}},
            {"data", bson.D{{"$push", "$$ROOT"}}}}}}
        projectStage := bson.D{
            {"$project", bson.D{
                {"_id", 0},
                {"total_count", 1},
                {"todo_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
            }},
        }

        result, err := todoCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing todo items"})
            return
        }
        defer result.Close(ctx)

        var todos []bson.M
        if err := result.All(ctx, &todos); err != nil {
            log.Fatal(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding todo items"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "data":    todos,
        })
    }
}



func CreateTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var todo models.Todo
        if err := c.BindJSON(&todo); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationErr := validate.Struct(todo)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }

        todo.Created_at = time.Now()
        todo.Updated_at = time.Now()
        todo.ID = primitive.NewObjectID()
        todo.User_uid=c.GetString("uid")

        _, err := todoCollection.InsertOne(ctx, todo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Todo created successfully",
            "data":    todo,
        })
    }
}




func UpdateTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        todoID := c.Param("todo_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var updatedTodo models.Todo
        if err := c.BindJSON(&updatedTodo); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationErr := validate.Struct(updatedTodo)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }

        objID, err := primitive.ObjectIDFromHex(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo_id format"})
            return
        }

        filter := bson.M{"_id": objID}
        update := bson.M{
            "$set": bson.M{
                "title":       updatedTodo.Title,
                "description": updatedTodo.Description,
                "status":      updatedTodo.Status,
                "updated_at":  time.Now(),
            },
        }

        _, err = todoCollection.UpdateOne(ctx, filter, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Todo updated successfully",
            "data":    updatedTodo,
        })
    }
}




func DeleteTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        todoID := c.Param("todo_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        objID, err := primitive.ObjectIDFromHex(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo_id format"})
            return
        }

        filter := bson.M{"_id": objID}

        _, err = todoCollection.DeleteOne(ctx, filter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Todo deleted successfully",
        })
    }
}
