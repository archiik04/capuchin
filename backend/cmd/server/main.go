package main

import (
	"capuchin/internal/auth"
	"capuchin/internal/database"
	"capuchin/internal/middleware"
	"capuchin/internal/models"
	"capuchin/internal/store"

	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// path to db
const dbPath = "../../db/db.json"

var todos []models.Todo
var mu sync.RWMutex

func main() {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Connect Database
	database.Connect()

	// Load todos once
	todos, err = store.Load(dbPath)
	if err != nil {
		todos = []models.Todo{}
	}

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// todos

	r.GET("/todos", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()
		c.JSON(200, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		var newTodo models.Todo

		if err := c.ShouldBindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newTodo.ID = uuid.New().String()

		mu.Lock()
		todos = append(todos, newTodo)
		store.Save(dbPath, todos)
		mu.Unlock()

		c.JSON(http.StatusOK, newTodo)
	})

	r.PATCH("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		mu.Lock()
		defer mu.Unlock()

		for i, t := range todos {
			if t.ID == id {
				todos[i].Completed = !todos[i].Completed
				store.Save(dbPath, todos)
				c.JSON(200, todos[i])
				return
			}
		}

		c.JSON(404, gin.H{"message": "Todo not found"})
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		mu.Lock()
		defer mu.Unlock()

		for i, t := range todos {
			if t.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				store.Save(dbPath, todos)
				c.JSON(200, gin.H{"message": "Todo deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"message": "Todo not found"})
	})

	r.PATCH("/todos/:id/edit", func(c *gin.Context) {
		id := c.Param("id")

		var body struct {
			Item string `json:"item"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for i, t := range todos {
			if t.ID == id {
				todos[i].Item = body.Item
				store.Save(dbPath, todos)
				c.JSON(200, todos[i])
				return
			}
		}

		c.JSON(404, gin.H{"message": "Todo not found"})
	})

	// Authentication
	r.POST("/signup", auth.SignupHandler)
	r.POST("/login", auth.LoginHandler)

	// Authenticated user todos (require JWT)
	userTodos := r.Group("/user")
	userTodos.Use(middleware.AuthRequired())
	{
		// Get all todos for authenticated user
		userTodos.GET("/todos", func(c *gin.Context) {
			userID := c.GetUint("userID")

			var todos []models.UserTodo
			result := database.DB.Where("user_id = ?", userID).Find(&todos)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
				return
			}

			c.JSON(http.StatusOK, todos)
		})

		// Create a new todo for authenticated user
		userTodos.POST("/todos", func(c *gin.Context) {
			userID := c.GetUint("userID")

			var input struct {
				Item      string `json:"item"`
				Completed bool   `json:"completed"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			todo := models.UserTodo{
				UserID:    userID,
				Item:      input.Item,
				Completed: input.Completed,
			}

			result := database.DB.Create(&todo)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
				return
			}

			c.JSON(http.StatusOK, todo)
		})

		// Toggle completion status
		userTodos.PATCH("/todos/:id", func(c *gin.Context) {
			userID := c.GetUint("userID")
			id := c.Param("id")

			var todo models.UserTodo
			result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo)

			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			todo.Completed = !todo.Completed
			database.DB.Save(&todo)

			c.JSON(http.StatusOK, todo)
		})

		// Delete a todo
		userTodos.DELETE("/todos/:id", func(c *gin.Context) {
			userID := c.GetUint("userID")
			id := c.Param("id")

			result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserTodo{})

			if result.Error != nil || result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
		})

		// Edit a todo's text
		userTodos.PATCH("/todos/:id/edit", func(c *gin.Context) {
			userID := c.GetUint("userID")
			id := c.Param("id")

			var input struct {
				Item string `json:"item"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			var todo models.UserTodo
			result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo)

			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			todo.Item = input.Item
			database.DB.Save(&todo)

			c.JSON(http.StatusOK, todo)
		})
	}

	r.Run(":8080")
}

func saveTodos() {
	data, _ := json.MarshalIndent(todos, "", "  ")
	os.WriteFile(dbPath, data, 0644)
}
