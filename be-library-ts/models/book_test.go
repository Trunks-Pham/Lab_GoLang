package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookCreation(t *testing.T) {
	book := Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Description: "Test Description",
		PublishedAt: "2024-10-16",
	}

	assert.NotNil(t, book.ID)
	assert.NotNil(t, book.CreatedAt)
	assert.NotNil(t, book.UpdatedAt)
}

func TestBookUpdate(t *testing.T) {
	book := Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Description: "Test Description",
		PublishedAt: "2024-10-16",
	}

	book.Title = "Updated Title"
	assert.NotNil(t, book.ID)
	assert.NotNil(t, book.CreatedAt)
	assert.NotNil(t, book.UpdatedAt)
}
