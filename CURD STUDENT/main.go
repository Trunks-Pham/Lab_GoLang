package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var students = make(map[int]string) // Sử dụng map để lưu trữ dữ liệu

var idCounter = 1

func main() {
	r := gin.Default()

	initData()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the Student API!")
	})

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, world")
	})

	r.GET("/students", GetStudents)
	r.POST("/students", CreateStudent)
	r.GET("/students/:id", GetStudent)
	r.PUT("/students/:id", UpdateStudent)
	r.DELETE("/students/:id", DeleteStudent)

	r.Run() // Mặc định chạy trên cổng 8080
}

// Hàm khởi tạo dữ liệu từ file JSON
func initData() {
	file, err := os.Open("students.json")
	if err != nil {
		// Nếu không thể mở file, khởi tạo dữ liệu mặc định
		students[1] = "John Doe"
		students[2] = "Jane Smith"
		students[3] = "Alice Johnson"
		students[4] = "Bob Brown"
		saveData() // Lưu dữ liệu mặc định vào file
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &students)
	if err != nil {
		panic(err)
	}

	// Cập nhật idCounter dựa trên dữ liệu hiện có
	for id := range students {
		if id >= idCounter {
			idCounter = id + 1
		}
	}
}

// Lưu dữ liệu vào file JSON
func saveData() {
	data, err := json.Marshal(students)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("students.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

// GetStudents trả về danh sách tất cả các sinh viên
func GetStudents(c *gin.Context) {
	c.JSON(http.StatusOK, students)
}

// CreateStudent tạo một sinh viên mới
func CreateStudent(c *gin.Context) {
	var student struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := idCounter
	idCounter++
	students[id] = student.Name
	saveData() // Lưu dữ liệu vào file
	c.JSON(http.StatusCreated, gin.H{"id": id, "name": student.Name})
}

// GetStudent trả về thông tin một sinh viên dựa trên ID
func GetStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	name, exists := students[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
}

// UpdateStudent cập nhật thông tin của một sinh viên dựa trên ID
func UpdateStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var student struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := students[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	students[id] = student.Name
	saveData() // Lưu dữ liệu vào file
	c.JSON(http.StatusOK, gin.H{"id": id, "name": student.Name})
}

// DeleteStudent xóa một sinh viên dựa trên ID
func DeleteStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if _, exists := students[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	delete(students, id)
	saveData() // Lưu dữ liệu vào file
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
}
