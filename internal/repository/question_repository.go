package repository

import (
	"context"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Define the interface for basic CRUD operations
type QuestionRepositoryInterface interface {
	CreateQuestion(question model.Question) (*model.Question, error) // Return the ID as a string
	GetQuestionByID(id primitive.ObjectID) (*model.Question, error)
	GetAllQuestions() ([]model.Question, error)
	UpdateQuestion(id primitive.ObjectID, question model.Question) (bool, error) // Return success status
	DeleteQuestion(id primitive.ObjectID) (bool, error)                          // Return success status
}

type QuestionRepository struct {
	collection *mongo.Collection
}

// NewQuestionRepository creates a new repository instance
func NewQuestionRepository(db *mongo.Database) *QuestionRepository {
	return &QuestionRepository{
		collection: db.Collection("questions"),
	}
}

// Implementing the interface methods

func (r *QuestionRepository) CreateQuestion(question model.Question) (*model.Question, error) {
	// Set the CreatedAt and UpdatedAt timestamps
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	// Insert the question into the database
	result, err := r.collection.InsertOne(context.Background(), question)
	if err != nil {
		return nil, err
	}

	// Set the InsertedID to the question's ID field (make sure it's a valid ObjectID)
	question.ID = result.InsertedID.(primitive.ObjectID)

	// Return the updated question with its ID
	return &question, nil
}

func (r *QuestionRepository) GetQuestionByID(id primitive.ObjectID) (*model.Question, error) {
	var question model.Question
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return the custom "not found" error
			return nil, customerrors.New(404, "Question not found for ID: "+id.Hex())
		}
		// Return other errors as-is
		return nil, err
	}
	return &question, err
}

func (r *QuestionRepository) GetAllQuestions() ([]model.Question, error) {
	var questions []model.Question
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var question model.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}

func (r *QuestionRepository) UpdateQuestion(id primitive.ObjectID, question model.Question) (bool, error) {
	question.UpdatedAt = time.Now()
	updateResult, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": question})
	if err != nil {
		return false, customerrors.ErrInternal
	}

	// Check if no document was matched
	if updateResult.MatchedCount == 0 {
		return false, customerrors.New(404, "Question not found with ID: "+id.Hex())
	}

	return true, nil
}

func (r *QuestionRepository) DeleteQuestion(id primitive.ObjectID) (bool, error) {
	deleteResult, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, customerrors.ErrInternal
	}

	// Check if no document was deleted
	if deleteResult.DeletedCount == 0 {
		return false, customerrors.New(404, "Question not found with ID: "+id.Hex())
	}

	return true, nil
}
