// internal/repository/question_repository.go
package repository

import (
	"context"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var questionCollection *mongo.Collection

func InitQuestionCollection() {
	if db.Client != nil {
		questionCollection = db.Client.Database("skillcode_db").Collection("questions")
	} else {
		log.Fatal("MongoDB client is not initialized")
	}
}

func CreateQuestion(question models.Question) (*mongo.InsertOneResult, error) {
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()
	return questionCollection.InsertOne(context.Background(), question)
}

func GetQuestionByID(id primitive.ObjectID) (models.Question, error) {
	var question models.Question
	err := questionCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
	return question, err
}

func GetAllQuestions() ([]models.Question, error) {
	var questions []models.Question
	cursor, err := questionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var question models.Question
		cursor.Decode(&question)
		questions = append(questions, question)
	}
	return questions, nil
}

func UpdateQuestion(id primitive.ObjectID, question models.Question) (*mongo.UpdateResult, error) {
	question.UpdatedAt = time.Now()
	return questionCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": question})
}

func DeleteQuestion(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return questionCollection.DeleteOne(context.Background(), bson.M{"_id": id})
}
