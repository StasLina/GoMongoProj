package repository

import (
	"LibrarySystem/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookInventoryRepository struct {
	collection *mongo.Collection
}

func NewBookInventoryRepository(db *mongo.Database) *BookInventoryRepository {
	return &BookInventoryRepository{
		collection: db.Collection("BookInventory"),
	}
}

func (r *BookInventoryRepository) GetAllBooksInventory(ctx context.Context) ([]models.BookInventory, error) {
	var books []models.BookInventory
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book models.BookInventory
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *BookInventoryRepository) GetBookInventoryByID(ctx context.Context, id primitive.ObjectID) (*models.BookInventory, error) {
	var book models.BookInventory
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	return &book, err
}

func (r *BookInventoryRepository) CreateBookInventory(ctx context.Context, book *models.BookInventory) error {
	_, err := r.collection.InsertOne(ctx, book)
	return err
}

func (r *BookInventoryRepository) UpdateBookInventory(ctx context.Context, book *models.BookInventory) error {
	filter := bson.M{"_id": book.ID}
	update := bson.M{
		"$set": bson.M{
			"bookId":     book.BookID,
			"locationId": book.LocationID,
			"title":      book.Title,
			"author":     book.Author,
			"quantity":   book.Quantity,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *BookInventoryRepository) DeleteBookInventory(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Генерация нового ObjectID
func (r *BookInventoryRepository) GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
