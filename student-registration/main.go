package main

import (
	"student-registration/models"
	"student-registration/routers"
)

func main() {
	models.InitDB()
	r := routers.SetupRouter()
	r.Run(":8080") // Chạy server ở port 8080
}
