package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HardBook struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InventoryID     primitive.ObjectID `bson:"inventoryId" json:"inventoryId"`
	AcquisitionDate primitive.DateTime `bson:"acquisitionDate" json:"acquisitionDate"`
}
