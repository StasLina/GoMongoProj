package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookItem struct {
	HardBookID primitive.ObjectID `bson:"hardBookId,omitempty"`
	Title      string             `bson:"title,omitempty"`
	Author     string             `bson:"author,omitempty"`
}

type BookIssueUser struct {
	BookIssuesID primitive.ObjectID `bson:"bookIssuesId,omitempty"`
	FullName     string             `bson:"fullName,omitempty"`
}

type BookIssue struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Book          BookItem           `bson:"book,omitempty"`
	BookIssueUser BookIssueUser      `bson:"bookIssueUser,omitempty"`
	IssueDate     time.Time          `bson:"issueDate,omitempty"`
	DueDate       time.Time          `bson:"dueDate,omitempty"`
	ReturnDate    *time.Time         `bson:"returnDate,omitempty"` // может быть nil
	LocationID    primitive.ObjectID `bson:"locationId,omitempty"`
	Status        string             `bson:"status,omitempty"` // "возвращено", "задолженность"
}
