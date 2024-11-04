package controllers

import (
	"net/http"

	"student-registration/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mã hóa mật khẩu trước khi lưu vào cơ sở dữ liệu
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(student.Password), 10)
	student.Password = string(hashedPassword)

	if err := models.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful!",
		"user": gin.H{
			"name":  student.Name,
			"email": student.Email,
		},
	})
}

func Login(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbStudent models.Student
	if err := models.DB.Where("email = ?", student.Email).First(&dbStudent).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbStudent.Password), []byte(student.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"user": gin.H{
			"id":    dbStudent.ID,
			"name":  dbStudent.Name,
			"email": dbStudent.Email,
		},
	})
}
