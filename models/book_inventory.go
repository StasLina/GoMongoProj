package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookInventory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID     primitive.ObjectID `bson:"bookId" json:"bookId"`
	LocationID primitive.ObjectID `bson:"locationId" json:"locationId"`
	Title      string             `bson:"title" json:"title"`
	Author     string             `bson:"author" json:"author"`
	Quantity   int                `bson:"quantity" json:"quantity"`
}
