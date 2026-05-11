package main

import (
	"task-tracker/database"
	"task-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	router := gin.Default()

	router.POST("/tasks", handlers.CreateTask)
	router.GET("/tasks", handlers.GetTasks)
	router.GET("/tasks/:id", handlers.GetTaskById)
	router.PATCH("/tasks/:id", handlers.UpdateTask)
	router.DELETE("/tasks/:id", handlers.DeleteTask)

	router.Run(":8080")
}