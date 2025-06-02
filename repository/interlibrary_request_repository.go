package repository

import (
	"LibrarySystem/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InterlibraryRequestRepository struct {
	collection *mongo.Collection
}

func NewInterlibraryRequestRepository(db *mongo.Database) *InterlibraryRequestRepository {
	return &InterlibraryRequestRepository{
		collection: db.Collection("InterlibraryRequests"),
	}
}

func (r *InterlibraryRequestRepository) GetAllInterlibraryRequests() ([]models.InterlibraryRequest, error) {
	var requests []models.InterlibraryRequest
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *InterlibraryRequestRepository) CreateInterlibraryRequest(req models.InterlibraryRequest) (*mongo.InsertOneResult, error) {
	req.ID = primitive.NewObjectID()
	return r.collection.InsertOne(context.TODO(), req)
}

func (r *InterlibraryRequestRepository) GetInterlibraryRequestByID(id primitive.ObjectID) (models.InterlibraryRequest, error) {
	var req models.InterlibraryRequest
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&req)
	return req, err
}

func (r *InterlibraryRequestRepository) UpdateInterlibraryRequest(id primitive.ObjectID, req models.InterlibraryRequest) (*mongo.UpdateResult, error) {
	return r.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": req})
}

func (r *InterlibraryRequestRepository) DeleteInterlibraryRequest(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
}

func (r *InterlibraryRequestRepository) GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
