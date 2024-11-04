package models

type Book struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Title         string `json:"title" bson:"title"`
	Genre         string `json:"genre" bson:"genre"`
	PublishedYear int    `json:"published_year" bson:"published_year"`
	AuthorID      string `json:"author_id" bson:"author_id"`
}
