package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"LibrarySystem/models"
)

type HardBookRepository struct {
	collection *mongo.Collection
}

func NewHardBookRepository(db *mongo.Database) *HardBookRepository {
	return &HardBookRepository{
		collection: db.Collection("HardBook"),
	}
}

func (r *HardBookRepository) GetAllHardBooks() ([]models.HardBook, error) {
	var books []models.HardBook
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var book models.HardBook
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *HardBookRepository) CreateHardBook(book *models.HardBook) error {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}
	if book.InventoryID.IsZero() {
		book.InventoryID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(context.TODO(), book)
	return err
}

func (r *HardBookRepository) GetHardBookByID(id string) (*models.HardBook, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	var book models.HardBook
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&book)
	return &book, err
}

func (r *HardBookRepository) UpdateHardBook(id string, updated *models.HardBook) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	updateData := bson.M{
		"$set": bson.M{
			"inventoryId":     updated.InventoryID,
			"acquisitionDate": updated.AcquisitionDate,
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, updateData)
	return err
}

func (r *HardBookRepository) DeleteHardBook(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	return err
}

func (r *HardBookRepository) GenerateObjectId() string {
	return primitive.NewObjectID().Hex()
}
