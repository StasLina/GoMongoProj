package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	LocationID primitive.ObjectID `bson:"locationId" json:"locationId"`
	Name       string             `bson:"name" json:"name"`
	Type       string             `bson:"type" json:"type"`
}

type Library struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Address   string             `bson:"address" json:"address"`
	Locations []Location         `bson:"locations" json:"locations"`
}

func ToHex(oid primitive.ObjectID) string {
	return oid.Hex()
}
