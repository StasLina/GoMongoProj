// repository/book_repository.go
package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"LibrarySystem/models"
)

type BookRepository struct {
	collection *mongo.Collection
}

func NewBookRepository(db *mongo.Database) *BookRepository {
	return &BookRepository{
		collection: db.Collection("Books"),
	}
}

func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) CreateBook(book *models.Book) error {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(context.Background(), book)
	return err
}

func (r *BookRepository) GetBookByID(id primitive.ObjectID) (*models.Book, error) {
	var book models.Book
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book)
	return &book, err
}

func (r *BookRepository) UpdateBook(id primitive.ObjectID, updated *models.Book) error {
	update := bson.M{
		"$set": bson.M{
			"title":           updated.Title,
			"author":          updated.Author,
			"publicationYear": updated.PublicationYear,
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}

func (r *BookRepository) DeleteBook(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *BookRepository) GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
