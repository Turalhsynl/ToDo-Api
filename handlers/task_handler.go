package handlers

import (
	"net/http"
	"strings"
	"task-tracker/database"
	"task-tracker/models"

	"github.com/gin-gonic/gin"
	"time"
)

func CreateTask(c *gin.Context) {

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})

		return
	}

	task.Title = strings.TrimSpace(task.Title)

	if task.Title == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})

		return
	}

	if len(task.Title) > 200 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title max length is 200",
		})

		return
	}

	if task.Status != "pending" && task.Status != "done" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Status must be pending or done",
		})

		return
	}

	database.DB.Create(&task)

	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {

	var tasks []models.Task

	status := c.Query("status")

	if status != "" {

		database.DB.
			Where("status = ?", status).
			Order("created_at desc").
			Find(&tasks)

	} else {

		database.DB.
			Order("created_at desc").
			Find(&tasks)
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTaskById(c *gin.Context) {

	id := c.Param("id")

	var task models.Task

	result := database.DB.First(&task, id)

	if result.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})

		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {

	id := c.Param("id")

	var task models.Task

	result := database.DB.First(&task, id)

	if result.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})

		return
	}

	var updatedTask map[string]interface{}

	if err := c.ShouldBindJSON(&updatedTask); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})

		return
	}

	if title, exists := updatedTask["title"]; exists {

		titleStr := strings.TrimSpace(title.(string))

		if titleStr == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Title is required",
			})

			return
		}

		if len(titleStr) > 200 {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Title max length is 200",
			})

			return
		}

		task.Title = titleStr
	}

	if status, exists := updatedTask["status"]; exists {

		statusStr := status.(string)

		if statusStr != "pending" && statusStr != "done" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status must be pending or done",
			})

			return
		}

		if statusStr == "done" && strings.TrimSpace(task.Title) == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Task cannot be marked done without title",
			})

			return
		}

		task.Status = statusStr
	}

	if dueDate, exists := updatedTask["due_date"]; exists {

	dateStr := dueDate.(string)

	parsedDate, err := time.Parse("2006-01-02", dateStr)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid due_date format. Use YYYY-MM-DD",
		})

		return
	}

	task.DueDate = &parsedDate
}

	database.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {

	id := c.Param("id")

	var task models.Task

	result := database.DB.First(&task, id)

	if result.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})

		return
	}

	database.DB.Delete(&task)

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}