// package routers

// import (
// 	"student-registration/controllers"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(cors.Default())

// 	r.POST("/register", controllers.Register)
// 	r.POST("/login", controllers.Login)
// 	r.GET("/courses", controllers.GetCourses)
// 	r.POST("/enroll", controllers.EnrollCourse)

// 	return r
// }

package routers

import (
	"student-registration/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// Đăng ký và đăng nhập
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/students", controllers.GetAllStudents)
	r.DELETE("/students/:id", controllers.DeleteStudent)
	r.PUT("/students/:id", controllers.UpdateInfor_Student)

	// Các endpoint cho khóa học
	r.GET("/courses", controllers.GetCourses)          // Lấy danh sách khóa học
	r.GET("/courses/:id", controllers.GetCourseByID)   // Lấy thông tin khóa học theo ID
	r.POST("/courses", controllers.CreateCourse)       // Tạo khóa học mới
	r.PUT("/courses/:id", controllers.UpdateCourse)    // Cập nhật thông tin khóa học
	r.DELETE("/courses/:id", controllers.DeleteCourse) // Xóa khóa học

	// Đăng ký khóa học
	r.POST("/enroll", controllers.EnrollCourse) // Đăng ký khóa học

	// Các endpoint cho sinh viên
	r.GET("/students/:id", controllers.GetStudentByID)                        // Lấy thông tin sinh viên theo ID
	r.GET("/students/:id/courses", controllers.GetEnrolledCoursesByStudentID) // Lấy danh sách khóa học của sinh viên

	return r
}
