package main

import (
	"library_management/controllers"
	"library_management/database"
	"library_management/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Kết nối Database
	database.Connect()
	database.DB.AutoMigrate(&models.Book{})

	r := gin.Default()

	// Các route API
	r.GET("/books", controllers.GetBooks)
	r.POST("/books", controllers.CreateBook)
	r.PUT("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	r.Run(":8080")
}
