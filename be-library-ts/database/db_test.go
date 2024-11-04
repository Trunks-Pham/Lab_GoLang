package database

import (
	"library_management/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	Connect()
	if DB == nil {
		t.Errorf("Database connection not established")
	}
}

func TestConnectError(t *testing.T) {
	// Thay đổi dsn để tạo lỗi kết nối
	dsn := "postgresql://wrong_user:wrong_password@wrong_host/wrong_db"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestDB(t *testing.T) {
	Connect()
	if DB == nil {
		t.Errorf("Database connection not established")
	}
	// Test các phương thức của DB
	var count int64
	DB.Model(&models.Book{}).Count(&count)
	if count < 0 {
		t.Errorf("Expected count >= 0, but got %d", count)
	}
}
