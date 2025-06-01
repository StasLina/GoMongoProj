package repository

import (
	"LibrarySystem/db"
	"LibrarySystem/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Library = models.Library

func GetAllLibraries() ([]Library, error) {
	var libraries []Library

	collection := db.GetCollection("Libraries")

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &libraries); err != nil {
		return nil, err
	}

	fmt.Println("Количество библиотек:", len(libraries))

	return libraries, nil
}

func GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func CreateLibrary(library models.Library) error {
	collection := db.GetCollection("Libraries")

	_, err := collection.InsertOne(context.TODO(), library)
	return err
}

func GetLibraryByID(id primitive.ObjectID) (models.Library, error) {
	collection := db.GetCollection("Libraries")

	var library models.Library
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&library)
	if err != nil {
		return models.Library{}, err
	}

	return library, nil
}

func UpdateLibrary(id primitive.ObjectID, library models.Library) error {
	collection := db.GetCollection("Libraries")

	update := bson.M{
		"$set": bson.M{
			"name":      library.Name,
			"address":   library.Address,
			"locations": library.Locations,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func DeleteLibrary(id primitive.ObjectID) error {
	collection := db.GetCollection("Libraries")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
