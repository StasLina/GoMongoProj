package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterlibraryRequest struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	SourceLocation      LocationItem       `bson:"sourceLocation"`
	DestinationLocation LocationItem       `bson:"destinationLocation"`
	RequestDate         primitive.DateTime `bson:"requestDate"`
	Status              string             `bson:"status"`
	Details             []BookDetail       `bson:"details"`
}

type LocationItem struct {
	LocationID primitive.ObjectID `bson:"locationId"`
	LibraryID  primitive.ObjectID `bson:"libraryId"`
}

type BookDetail struct {
	BookID   primitive.ObjectID `bson:"bookId"`
	Title    string             `bson:"title"`
	Author   string             `bson:"author"`
	Quantity int                `bson:"quantity"`
}
