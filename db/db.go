package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	uri := "mongodb://root:example@192.192.56.2:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	Client = client
	fmt.Println("✅ Подключено к MongoDB!")
}

func GetCollection(collecionName string) *mongo.Collection {
	return Client.Database("LibrarySystem").Collection(collecionName)
}
