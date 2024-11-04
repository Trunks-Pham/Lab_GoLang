// package controllers

// import (
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"student-registration/models"

// 	"github.com/gin-gonic/gin"
// )

// func EnrollCourse(c *gin.Context) {
// 	var enrollment models.Enrollment
// 	if err := c.ShouldBindJSON(&enrollment); err != nil {
// 		log.Printf("Error binding JSON: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
// 		return
// 	}

// 	// Kiểm tra xem khóa học có tồn tại không
// 	var course models.Course
// 	if err := models.DB.First(&course, enrollment.CourseID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
// 		return
// 	}

// 	if err := models.DB.Create(&enrollment).Error; err != nil {
// 		log.Printf("Error enrolling in course: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Enrollment successful"})
// }

// func GetCourses(c *gin.Context) {
// 	var courses []models.Course
// 	if err := models.DB.Find(&courses).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, courses)
// }

// // Lấy thông tin sinh viên theo ID
// func GetStudentByID(c *gin.Context) {
// 	studentID := c.Param("id")
// 	id, err := strconv.Atoi(studentID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
// 		return
// 	}

// 	var student models.Student
// 	if err := models.DB.First(&student, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, student)
// }

// // Lấy thông tin khóa học theo ID
// func GetCourseByID(c *gin.Context) {
// 	courseID := c.Param("id")
// 	id, err := strconv.Atoi(courseID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
// 		return
// 	}

// 	var course models.Course
// 	if err := models.DB.First(&course, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, course)
// }

// // Lấy danh sách các khóa học mà sinh viên đã đăng ký
// func GetEnrolledCoursesByStudentID(c *gin.Context) {
// 	studentID := c.Param("id")
// 	id, err := strconv.Atoi(studentID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
// 		return
// 	}

// 	var courses []models.Course
// 	if err := models.DB.Table("courses").Select("courses.id, courses.name, courses.description").
// 		Joins("join enrollments on courses.id = enrollments.course_id").
// 		Where("enrollments.student_id = ?", id).Scan(&courses).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, courses)
// }

package controllers

import (
	"log"
	"net/http"
	"strconv"

	"student-registration/models"

	"github.com/gin-gonic/gin"
)

// EnrollCourse - Đăng ký khóa học
func EnrollCourse(c *gin.Context) {
	var enrollment models.Enrollment
	if err := c.ShouldBindJSON(&enrollment); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Kiểm tra xem khóa học có tồn tại không
	var course models.Course
	if err := models.DB.First(&course, enrollment.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if err := models.DB.Create(&enrollment).Error; err != nil {
		log.Printf("Error enrolling in course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment successful"})
}

// GetCourses - Lấy danh sách tất cả các khóa học
func GetCourses(c *gin.Context) {
	var courses []models.Course
	if err := models.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// GetStudentByID - Lấy thông tin sinh viên theo ID
func GetStudentByID(c *gin.Context) {
	studentID := c.Param("id")
	id, err := strconv.Atoi(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := models.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// UpdateInfor_Student
func UpdateInfor_Student(c *gin.Context) {
	studentID := c.Param("id")
	id, err := strconv.Atoi(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := models.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := models.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student information"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// GetAllStudents
func GetAllStudents(c *gin.Context) {
	var students []models.Student
	if err := models.DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// DeleteStudent
func DeleteStudent(c *gin.Context) {
	studentID := c.Param("id")
	id, err := strconv.Atoi(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := models.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if err := models.DB.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

// GetCourseByID - Lấy thông tin khóa học theo ID
func GetCourseByID(c *gin.Context) {
	courseID := c.Param("id")
	id, err := strconv.Atoi(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := models.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// GetEnrolledCoursesByStudentID - Lấy danh sách các khóa học mà sinh viên đã đăng ký
func GetEnrolledCoursesByStudentID(c *gin.Context) {
	studentID := c.Param("id")
	id, err := strconv.Atoi(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var courses []models.Course
	if err := models.DB.Table("courses").Select("courses.id, courses.name, courses.description").
		Joins("join enrollments on courses.id = enrollments.course_id").
		Where("enrollments.student_id = ?", id).Scan(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// UpdateCourse - Cập nhật thông tin khóa học
func UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, err := strconv.Atoi(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course.ID = uint(id)
	if err := models.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// CreateCourse - Tạo khóa học mới
func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := models.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully", "course": course})
}

// DeleteCourse - Xóa khóa học
func DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")
	id, err := strconv.Atoi(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := models.DB.Delete(&models.Course{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
