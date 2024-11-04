package controllers

import (
	"bytes"
	"encoding/json"
	"library_management/database"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Hàm setupRouter để cấu hình router cho các test case
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", GetBooks)
	r.POST("/books", CreateBook)
	r.PUT("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)
	return r
}

// Hàm để reset cơ sở dữ liệu giữa các test
func setupTestDB() {
	database.Connect()
	database.DB.Exec("DROP TABLE IF EXISTS books") // Xóa bảng để tránh dữ liệu thừa
	database.DB.AutoMigrate(&models.Book{})        // Tạo lại bảng
}

// Test tạo một cuốn sách hợp lệ
func TestCreateBook(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	book := models.Book{Title: "Test Book", Author: "Test Author", Description: "Test Description", PublishedAt: "2024-10-16"}
	jsonBook, _ := json.Marshal(book)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBook))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdBook models.Book
	json.Unmarshal(w.Body.Bytes(), &createdBook)

	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author, createdBook.Author)
}

// Test tạo sách với dữ liệu không hợp lệ (thiếu tiêu đề)
func TestCreateBookInvalid(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Test với sách thiếu tiêu đề
	book := models.Book{Title: "", Author: "Test Author", Description: "Test Description", PublishedAt: "2024-10-16"}
	jsonBook, _ := json.Marshal(book)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBook))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test lấy danh sách sách
func TestGetBooks(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Thêm một cuốn sách để có dữ liệu kiểm tra
	book := models.Book{Title: "Get Test Book", Author: "Test Author", Description: "Test Description", PublishedAt: "2024-10-16"}
	database.DB.Create(&book)

	req, _ := http.NewRequest("GET", "/books?title=Get", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var books []models.Book
	json.Unmarshal(w.Body.Bytes(), &books)
	assert.Len(t, books, 1)
	assert.Equal(t, book.Title, books[0].Title)
}

// Test lấy danh sách sách nhưng không có sách nào
func TestGetBooksEmpty(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var books []models.Book
	json.Unmarshal(w.Body.Bytes(), &books)
	assert.Len(t, books, 0)
}

// Test cập nhật thông tin một cuốn sách
func TestUpdateBook(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Tạo một cuốn sách để cập nhật
	book := models.Book{Title: "Update Test Book", Author: "Test Author", Description: "Updated Description", PublishedAt: "2024-10-16"}
	database.DB.Create(&book)

	book.Title = "Updated Book"
	jsonBook, _ := json.Marshal(book)

	req, _ := http.NewRequest("PUT", "/books/"+strconv.FormatUint(uint64(book.ID), 10), bytes.NewBuffer(jsonBook))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedBook models.Book
	database.DB.First(&updatedBook, book.ID)
	assert.Equal(t, "Updated Book", updatedBook.Title)
}

// Test cập nhật sách với ID không tồn tại
func TestUpdateBookNotFound(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Sách không tồn tại
	book := models.Book{Title: "Non-existent Book", Author: "Test Author", Description: "Updated Description", PublishedAt: "2024-10-16"}
	book.ID = 999 // ID không tồn tại
	jsonBook, _ := json.Marshal(book)

	req, _ := http.NewRequest("PUT", "/books/999", bytes.NewBuffer(jsonBook))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test xóa một cuốn sách
func TestDeleteBook(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Tạo một cuốn sách để xóa
	book := models.Book{Title: "Delete Test Book", Author: "Test Author", Description: "Delete Description", PublishedAt: "2024-10-16"}
	database.DB.Create(&book)

	req, _ := http.NewRequest("DELETE", "/books/"+strconv.FormatUint(uint64(book.ID), 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	var deletedBook models.Book
	result := database.DB.First(&deletedBook, book.ID)
	assert.Error(t, result.Error) // Đảm bảo sách đã bị xóa
}

// Test xóa sách không tồn tại
func TestDeleteBookNotFound(t *testing.T) {
	setupTestDB()

	r := setupRouter()

	// Sách không tồn tại
	req, _ := http.NewRequest("DELETE", "/books/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
