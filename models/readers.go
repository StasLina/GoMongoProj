package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reader struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName      string             `bson:"fullName" json:"fullName"`
	Category      Category           `bson:"category" json:"category"`
	Attributes    []Attribute        `bson:"attributes" json:"attributes"`
	Subscriptions []Subscription     `bson:"subscriptions" json:"subscriptions"`
	Fines         []Fine             `bson:"fines" json:"fines"`
}

type Category struct {
	CategoryID    primitive.ObjectID `bson:"categoryId" json:"categoryId"`
	Name          string             `bson:"name" json:"name"`
	MaxBorrowDays int                `bson:"maxBorrowDays" json:"maxBorrowDays"`
}

type Attribute struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

type Subscription struct {
	LocationID primitive.ObjectID `bson:"locationId" json:"locationId"`
	StartDate  primitive.DateTime `bson:"startDate" json:"startDate"`
	EndDate    primitive.DateTime `bson:"endDate" json:"endDate"`
	MaxBooks   int                `bson:"maxBooks" json:"maxBooks"`
}

type Fine struct {
	Amount      float64            `bson:"amount" json:"amount"`
	Reason      string             `bson:"reason" json:"reason"`
	FineDate    primitive.DateTime `bson:"fineDate" json:"fineDate"`
	FineDateEnd primitive.DateTime `bson:"fineDateEnd,omitempty" json:"fineDateEnd"`
}
