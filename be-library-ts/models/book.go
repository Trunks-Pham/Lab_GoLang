package models

import "gorm.io/gorm"

// Cấu trúc Book
type Book struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
}
