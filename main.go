package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Todo represents a simple todo item
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

var todos []Todo

func main() {
	r := gin.Default()

	// Endpoint to get all todos
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})

	// Endpoint to get a specific todo by ID
	r.GET("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, todo := range todos {
			if todo.ID == id {
				c.JSON(http.StatusOK, todo)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	})

	// Endpoint to create a new todo
	r.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Clean the ID field
		newTodo.ID = cleanID(newTodo.ID)

		// Generate a unique ID for the new todo
		newTodo.ID = generateID()

		// Add the new todo to the list
		todos = append(todos, newTodo)

		c.JSON(http.StatusCreated, newTodo)
	})

	// Endpoint to update an existing todo by ID
	r.PUT("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedTodo Todo
		if err := c.BindJSON(&updatedTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Find the todo by ID
		for i, todo := range todos {
			if todo.ID == id {
				// Update the todo
				todos[i] = updatedTodo
				c.JSON(http.StatusOK, updatedTodo)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	})

	// Endpoint to delete a todo by ID
	r.DELETE("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, todo := range todos {
			if todo.ID == id {
				// Remove the todo from the list
				todos = append(todos[:i], todos[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	})

	// Run the server
	r.Run(":8080")
}

func cleanID(id string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, id)
}

func generateID() string {
	// This is a simple method to generate a unique ID. You might want to use a more robust method in production.
	return "id" + string(len(todos)+1)
}
