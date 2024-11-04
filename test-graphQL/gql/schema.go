package gql

import (
	"context"
	"project/db"
	"project/models"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Định nghĩa các GraphQL Object Types
var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.String},
			"name": &graphql.Field{Type: graphql.String},
			"age":  &graphql.Field{Type: graphql.Int},
			"booksWithGenres": &graphql.Field{
				Type: graphql.NewList(bookType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					author := params.Source.(models.Author)
					var books []models.Book
					collection := db.Client.Database("library").Collection("books")
					cursor, err := collection.Find(context.Background(), primitive.M{"authorID": author.ID})
					if err != nil {
						return nil, err
					}
					defer cursor.Close(context.Background())
					for cursor.Next(context.Background()) {
						var book models.Book
						if err := cursor.Decode(&book); err != nil {
							return nil, err
						}
						books = append(books, book)
					}
					return books, nil
				},
			},
		},
	},
)

var bookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"id":            &graphql.Field{Type: graphql.String},
			"title":         &graphql.Field{Type: graphql.String},
			"genre":         &graphql.Field{Type: graphql.String},
			"publishedYear": &graphql.Field{Type: graphql.Int},
			"authorID":      &graphql.Field{Type: graphql.String},
			"authorName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					book := params.Source.(models.Book)
					var author models.Author
					collection := db.Client.Database("library").Collection("authors")
					objID, _ := primitive.ObjectIDFromHex(book.AuthorID)
					err := collection.FindOne(context.Background(), primitive.M{"_id": objID}).Decode(&author)
					if err != nil {
						return nil, err
					}
					return author.Name, nil
				},
			},
		},
	},
)

// Root query cho GraphQL để lấy dữ liệu
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		// Lấy tất cả Authors
		"authors": &graphql.Field{
			Type: graphql.NewList(authorType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var authors []models.Author
				collection := db.Client.Database("library").Collection("authors")
				cursor, err := collection.Find(context.Background(), primitive.D{})
				if err != nil {
					return nil, err
				}
				defer cursor.Close(context.Background())
				for cursor.Next(context.Background()) {
					var author models.Author
					if err := cursor.Decode(&author); err != nil {
						return nil, err
					}
					authors = append(authors, author)
				}
				return authors, nil
			},
		},
		// Lấy tất cả Books
		"books": &graphql.Field{
			Type: graphql.NewList(bookType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var books []models.Book
				collection := db.Client.Database("library").Collection("books")
				cursor, err := collection.Find(context.Background(), primitive.D{})
				if err != nil {
					return nil, err
				}
				defer cursor.Close(context.Background())
				for cursor.Next(context.Background()) {
					var book models.Book
					if err := cursor.Decode(&book); err != nil {
						return nil, err
					}
					books = append(books, book)
				}
				return books, nil
			},
		},
		// Lấy Book theo ID
		"book": &graphql.Field{
			Type: bookType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				var book models.Book
				collection := db.Client.Database("library").Collection("books")
				objID, _ := primitive.ObjectIDFromHex(id)
				err := collection.FindOne(context.Background(), primitive.M{"_id": objID}).Decode(&book)
				if err != nil {
					return nil, err
				}
				return book, nil
			},
		},

		// Xóa Author
		"deleteAuthor": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				collection := db.Client.Database("library").Collection("authors")
				objID, _ := primitive.ObjectIDFromHex(id)

				// Xóa tất cả sách của tác giả
				bookCollection := db.Client.Database("library").Collection("books")
				_, err := bookCollection.DeleteMany(context.Background(), primitive.M{"authorID": id})
				if err != nil {
					return nil, err
				}

				// Xóa tác giả
				_, err = collection.DeleteOne(context.Background(), primitive.M{"_id": objID})
				if err != nil {
					return nil, err
				}

				return map[string]string{"id": id}, nil
			},
		},
	},
})

// Root mutation để tạo, cập nhật, xóa dữ liệu
var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// Tạo mới Author
		"createAuthor": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"age":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name := params.Args["name"].(string)
				age := params.Args["age"].(int)
				author := models.Author{Name: name, Age: age}
				collection := db.Client.Database("library").Collection("authors")
				result, err := collection.InsertOne(context.Background(), author)
				if err != nil {
					return nil, err
				}
				author.ID = result.InsertedID.(primitive.ObjectID).Hex()
				return author, nil
			},
		},
		// Cập nhật Author
		"updateAuthor": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"name": &graphql.ArgumentConfig{Type: graphql.String},
				"age":  &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				updateFields := primitive.M{}

				// Lấy các trường để cập nhật nếu có
				if name, ok := params.Args["name"].(string); ok {
					updateFields["name"] = name
				}
				if age, ok := params.Args["age"].(int); ok {
					updateFields["age"] = age
				}

				// Kết nối tới MongoDB
				collection := db.Client.Database("library").Collection("authors")
				objID, _ := primitive.ObjectIDFromHex(id)

				// Cập nhật thông tin tác giả
				_, err := collection.UpdateOne(context.Background(), primitive.M{"_id": objID}, primitive.M{"$set": updateFields})
				if err != nil {
					return nil, err
				}

				// Lấy lại tác giả đã cập nhật
				var updatedAuthor models.Author
				err = collection.FindOne(context.Background(), primitive.M{"_id": objID}).Decode(&updatedAuthor)
				if err != nil {
					return nil, err
				}

				return updatedAuthor, nil
			},
		},
		// Tạo mới Book
		"createBook": &graphql.Field{
			Type: bookType,
			Args: graphql.FieldConfigArgument{
				"title":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"genre":         &graphql.ArgumentConfig{Type: graphql.String},
				"publishedYear": &graphql.ArgumentConfig{Type: graphql.Int},
				"authorID":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				title := params.Args["title"].(string)
				genre, _ := params.Args["genre"].(string)
				publishedYear, _ := params.Args["publishedYear"].(int)
				authorID := params.Args["authorID"].(string)
				book := models.Book{Title: title, Genre: genre, PublishedYear: publishedYear, AuthorID: authorID}
				collection := db.Client.Database("library").Collection("books")
				result, err := collection.InsertOne(context.Background(), book)
				if err != nil {
					return nil, err
				}
				book.ID = result.InsertedID.(primitive.ObjectID).Hex()
				return book, nil
			},
		},
		// Cập nhật Book
		"updateBook": &graphql.Field{
			Type: bookType,
			Args: graphql.FieldConfigArgument{
				"id":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"title":         &graphql.ArgumentConfig{Type: graphql.String},
				"genre":         &graphql.ArgumentConfig{Type: graphql.String},
				"publishedYear": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				updateFields := primitive.M{}
				if title, ok := params.Args["title"].(string); ok {
					updateFields["title"] = title
				}
				if genre, ok := params.Args["genre"].(string); ok {
					updateFields["genre"] = genre
				}
				if publishedYear, ok := params.Args["publishedYear"].(int); ok {
					updateFields["publishedYear"] = publishedYear
				}
				collection := db.Client.Database("library").Collection("books")
				objID, _ := primitive.ObjectIDFromHex(id)
				_, err := collection.UpdateOne(context.Background(), primitive.M{"_id": objID}, primitive.M{"$set": updateFields})
				if err != nil {
					return nil, err
				}
				var updatedBook models.Book
				collection.FindOne(context.Background(), primitive.M{"_id": objID}).Decode(&updatedBook)
				return updatedBook, nil
			},
		},
		// Xóa Book
		"deleteBook": &graphql.Field{
			Type: bookType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				collection := db.Client.Database("library").Collection("books")
				objID, _ := primitive.ObjectIDFromHex(id)
				_, err := collection.DeleteOne(context.Background(), primitive.M{"_id": objID})
				if err != nil {
					return nil, err
				}
				return map[string]string{"id": id}, nil
			},
		},
		// Xóa Author
		"deleteAuthor": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				// Chuyển đổi ID sang ObjectID
				objID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return nil, err
				}

				// Kết nối tới MongoDB
				collection := db.Client.Database("library").Collection("authors")
				bookCollection := db.Client.Database("library").Collection("books")

				// Xóa tất cả sách có `authorID` là objID
				_, err = bookCollection.DeleteMany(context.Background(), primitive.M{"authorID": objID})
				if err != nil {
					return nil, err
				}

				// Xóa tác giả
				_, err = collection.DeleteOne(context.Background(), primitive.M{"_id": objID})
				if err != nil {
					return nil, err
				}

				return map[string]string{"id": id}, nil
			},
		},
	},
})

// Tạo schema GraphQL
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})
