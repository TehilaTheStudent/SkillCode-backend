package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/handlers"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Client
var questionHandler *handlers.QuestionHandler

// Setup integration test environment
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error
	testDB, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// Clean up database
	defer func() {
		testDB.Database("test_database").Drop(ctx)
		testDB.Disconnect(ctx)
	}()

	// Create repository, service, and handler
	questionRepo := repository.NewQuestionRepository(testDB.Database("test_database"))
	questionService := service.NewQuestionService(questionRepo)
	questionHandler = handlers.NewQuestionHandler(questionService)

	// Run tests
	m.Run()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	handlers.RegisterQuestionRoutes(router, questionHandler)
	return router
}

func insertTestQuestion(question model.Question) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := testDB.Database("test_database").Collection("questions")
	_, _ = collection.InsertOne(ctx, question)
}
func clearDatabase() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    testDB.Database("test_database").Drop(ctx)
}


func TestGetQuestionByID_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := model.Question{
		ID:          questionID,
		Title:       "Two Sum",
		Description: "Find two numbers that add up to a target.",
		Visibility:  "public",
		CreatedBy:   "user123",
	}
	insertTestQuestion(testQuestion)

	// Set up router
	router := setupRouter()

	// Act: Make the HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/"+questionID.Hex(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testQuestion.ID, response.ID)
	assert.Equal(t, testQuestion.Title, response.Title)
}
func TestCreateQuestion_Integration(t *testing.T) {
	clearDatabase()

	router := setupRouter()

	// Request body for the new question
	jsonBody := `{
        "title": "Binary Search",
        "description": "Implement a binary search algorithm.",
        "function_signature": "func binarySearch(nums []int, target int) int",
        "test_cases": [],
        "visibility": "public",
        "created_by": "user123"
    }`

	// Act: Make the HTTP POST request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response to verify the created question
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Binary Search", response["title"])
	assert.NotEmpty(t, response["id"])

	// Verify the question is in the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := testDB.Database("test_database").Collection("questions")
	var dbQuestion map[string]interface{}
	err := collection.FindOne(ctx, map[string]interface{}{"title": "Binary Search"}).Decode(&dbQuestion)
	assert.NoError(t, err)
	assert.Equal(t, "Binary Search", dbQuestion["title"])
}
func TestUpdateQuestion_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := model.Question{
		ID:                questionID,
		Title:             "Binary Search",
		Description:       "Implement a binary search algorithm.",
		FunctionSignature: "func binarySearch(nums []int, target int) int",
		Visibility:        "public",
		CreatedBy:         "user123",
	}
	insertTestQuestion(testQuestion)

	// Updated data
	updatedTitle := "Updated Binary Search"
	updatedDescription := "Implement an optimized binary search algorithm."
	jsonBody := fmt.Sprintf(`{
        "title": "%s",
        "description": "%s",
        "function_signature": "func binarySearch(nums []int, target int) int",
        "test_cases": [],
        "visibility": "public",
        "created_by": "user123"
    }`, updatedTitle, updatedDescription)

	router := setupRouter()

	// Act: Make the HTTP PUT request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/questions/"+questionID.Hex(), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, updatedTitle, response["title"])
	assert.Equal(t, updatedDescription, response["description"])

	// Verify the updated data in the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := testDB.Database("test_database").Collection("questions")
	var dbQuestion map[string]interface{}
	err := collection.FindOne(ctx, map[string]interface{}{"_id": questionID}).Decode(&dbQuestion)
	assert.NoError(t, err)
	assert.Equal(t, updatedTitle, dbQuestion["title"])
	assert.Equal(t, updatedDescription, dbQuestion["description"])
}


func TestDeleteQuestion_Integration(t *testing.T) {
	clearDatabase()
	
	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := model.Question{
		ID:          questionID,
		Title:       "Binary Search",
		Description: "Implement a binary search algorithm.",
		Visibility:  "public",
		CreatedBy:   "user123",
	}
	insertTestQuestion(testQuestion)

	router := setupRouter()

	// Act: Make the HTTP DELETE request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/questions/"+questionID.Hex(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the question is deleted from the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := testDB.Database("test_database").Collection("questions")
	var dbQuestion map[string]interface{}
	err := collection.FindOne(ctx, map[string]interface{}{"_id": questionID}).Decode(&dbQuestion)
	assert.Error(t, err) // Expect error because the question no longer exists
}
func TestTestQuestion_Integration(t *testing.T) {
	clearDatabase()

	router := setupRouter()

	// Arrange
	questionID := primitive.NewObjectID()
	testQuestion := model.Question{
		ID:          questionID,
		Title:       "Binary Search",
		Description: "Implement a binary search algorithm.",
		Visibility:  "public",
		CreatedBy:   "user123",
	}
	insertTestQuestion(testQuestion)

	userFunction := `func binarySearch(nums []int, target int) int {return 0 // Example implementation}`
	jsonBody := fmt.Sprintf(`{"user_function": "%s"}`, userFunction)

	// Act: Make the HTTP POST request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions/"+questionID.Hex()+"/test", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["result"])
}

func TestGetAllQuestions_Integration(t *testing.T) {
	clearDatabase()
	// Arrange: Insert multiple test questions
	questions := []model.Question{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Binary Search",
			Description: "Implement a binary search algorithm.",
			Visibility:  "public",
			CreatedBy:   "user123",
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Reverse String",
			Description: "Reverse a given string.",
			Visibility:  "private",
			CreatedBy:   "user456",
		},
	}
	for _, question := range questions {
		insertTestQuestion(question)
	}

	router := setupRouter()

	// Act: Make the HTTP GET request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response []model.Question
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, len(questions))
}