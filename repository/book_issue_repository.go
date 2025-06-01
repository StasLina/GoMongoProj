package repository

import (
	"LibrarySystem/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookIssueRepository struct {
	Collection *mongo.Collection
}

func NewBookIssueRepository(db *mongo.Database) *BookIssueRepository {
	return &BookIssueRepository{
		Collection: db.Collection("BookIssues"),
	}
}

func (r *BookIssueRepository) GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (r *BookIssueRepository) GetAllBookIssues(ctx context.Context) ([]models.BookIssue, error) {
	var issues []models.BookIssue
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (r *BookIssueRepository) CreateBookIssue(ctx context.Context, issue models.BookIssue) error {
	_, err := r.Collection.InsertOne(ctx, issue)
	return err
}

func (r *BookIssueRepository) GetBookIssueByID(ctx context.Context, id primitive.ObjectID) (models.BookIssue, error) {
	var issue models.BookIssue
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue)
	return issue, err
}

func (r *BookIssueRepository) UpdateBookIssue(ctx context.Context, id primitive.ObjectID, updated models.BookIssue) error {
	update := bson.M{
		"$set": bson.M{
			"book.hardBookId":            updated.Book.HardBookID,
			"book.title":                 updated.Book.Title,
			"book.author":                updated.Book.Author,
			"bookIssueUser.bookIssuesId": updated.BookIssueUser.BookIssuesID,
			"bookIssueUser.fullName":     updated.BookIssueUser.FullName,
			"issueDate":                  updated.IssueDate,
			"dueDate":                    updated.DueDate,
			"returnDate":                 updated.ReturnDate,
			"locationId":                 updated.LocationID,
			"status":                     updated.Status,
		},
	}
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *BookIssueRepository) DeleteBookIssue(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
