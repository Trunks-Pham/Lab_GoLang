package main

import (
	"project/db"
	"project/gql"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func main() {
	// Kết nối tới MongoDB
	db.Connect()

	r := gin.Default()

	// Khởi tạo handler GraphQL
	h := handler.New(&handler.Config{
		Schema:     &gql.Schema,
		Pretty:     true,
		GraphiQL:   true, // dùng GraphiQL để kiểm thử
		Playground: true,
	})

	// Tạo route cho GraphQL
	r.POST("/graphql", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	// Chạy server
	r.Run(":8080")
}