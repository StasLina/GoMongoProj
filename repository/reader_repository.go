package repository

import (
	"context"

	"LibrarySystem/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReaderRepository struct {
	collection *mongo.Collection
}

func NewReaderRepository(db *mongo.Database) *ReaderRepository {
	return &ReaderRepository{
		collection: db.Collection("Readers"),
	}
}

func (r *ReaderRepository) GetAllReaders() ([]models.Reader, error) {
	var readers []models.Reader
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var reader models.Reader
		if err := cursor.Decode(&reader); err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

func (r *ReaderRepository) GetReaderByID(id primitive.ObjectID) (*models.Reader, error) {
	var reader models.Reader
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&reader)
	if err != nil {
		return nil, err
	}
	return &reader, nil
}

func (r *ReaderRepository) CreateReader(reader *models.Reader) error {
	reader.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.Background(), reader)
	return err
}

func (r *ReaderRepository) UpdateReader(id primitive.ObjectID, reader *models.Reader) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": reader},
	)
	return err
}

func (r *ReaderRepository) DeleteReader(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *ReaderRepository) GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
