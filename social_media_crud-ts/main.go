package main

import (
	"social_media_crud/config"
	"social_media_crud/models"
	"social_media_crud/routes"
)

func main() {
	// Khởi tạo kết nối database và Redis
	config.InitDB()
	config.InitRedis()

	// Tạo bảng nếu chưa có
	config.DB.AutoMigrate(&models.User{}, &models.Post{})

	// Khởi tạo router
	r := routes.SetupRouter()
	r.Run(":8080")
}

// Chạy: docker-compose up -d
