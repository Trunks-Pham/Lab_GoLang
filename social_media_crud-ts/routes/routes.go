package routes

import (
	"social_media_crud/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.GET("/", controllers.GetUsers)
		userRoutes.GET("/:id", controllers.GetUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.PUT("/:id", controllers.UpdateUser)
	}

	postRoutes := r.Group("/posts")
	{
		postRoutes.POST("/", controllers.CreatePost)
		postRoutes.GET("/", controllers.GetPosts)
		postRoutes.GET("/:id", controllers.GetPost)
		postRoutes.DELETE("/:id", controllers.DeletePost)
		postRoutes.PUT("/:id", controllers.UpdatePost)
	}

	return r
}
