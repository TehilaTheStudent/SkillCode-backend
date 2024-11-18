package integration_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handler"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Client
var questionHandler *handler.QuestionHandler

// Setup integration test environment
func TestMain(m *testing.M) {
	utils.EnsureWorkingDirectory()
	setupDatabase()
	defer teardownDatabase()
	m.Run()
}

func setupDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load configuration
	cfg := config.LoadConfigAPI()

	// Connect to MongoDB
	var err error
	testDB, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		panic(err)
	}

	// Create repository, service, and handler
	questionRepo := repository.NewQuestionRepository(testDB.Database("test_database"))
	questionService := service.NewQuestionService(questionRepo)
	questionHandler = handler.NewQuestionHandler(questionService)
}

func teardownDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := testDB.Database("test_database").Drop(ctx)
	if err != nil {
		panic(err)
	}
	err = testDB.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	handler.RegisterQuestionRoutes(router, questionHandler)
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
	err := testDB.Database("test_database").Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func TestGetQuestionByID_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := utils.GenerateQuestion(map[string]interface{}{"ID": questionID})
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
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, testQuestion.ID, response.ID)
	assert.Equal(t, testQuestion.Title, response.Title)
}

func TestCreateQuestion_Integration(t *testing.T) {
	clearDatabase()

	// Setup the router
	router := setupRouter()

	// Generate request body for the new question
	jsonBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":       "Binary Search",
		"description": "Implement a binary search algorithm.",
		"test_cases": []map[string]string{
			{"input": "[1, 2, 3, 4, 5], 3", "expected_output": "2"},
			{"input": "[10, 20, 30, 40], 30", "expected_output": "2"},
		},
		"languages": []map[string]string{
			{"language": "golang", "function_signature": "func binarySearch(nums []int, target int) int"},
			{"language": "python", "function_signature": "def binary_search(nums, target):"},
		},
		"visibility": "public",
		"created_by": "user123",
	})

	// Act: Make the HTTP POST request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert: Check the HTTP response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response to verify the created question
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Binary Search", response["title"])
	assert.NotEmpty(t, response["id"])

	// Verify the question is stored in the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := testDB.Database("test_database").Collection("questions")
	var dbQuestion map[string]interface{}
	err = collection.FindOne(ctx, map[string]interface{}{"title": "Binary Search"}).Decode(&dbQuestion)

	assert.NoError(t, err)
	assert.Equal(t, "Binary Search", dbQuestion["title"])
	assert.Equal(t, "public", dbQuestion["visibility"])
	assert.Equal(t, "user123", dbQuestion["created_by"])
	assert.Len(t, dbQuestion["test_cases"], 2)
	assert.Len(t, dbQuestion["languages"], 2)
}

func TestUpdateQuestion_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := utils.GenerateQuestion(map[string]interface{}{"ID": questionID})
	insertTestQuestion(testQuestion)

	// Updated data
	updatedTitle := "Updated Binary Search"
	updatedDescription := "Implement an optimized binary search algorithm."
	jsonBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":       updatedTitle,
		"description": updatedDescription,
	})

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
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, updatedTitle, response["title"])
	assert.Equal(t, updatedDescription, response["description"])

	// Verify the updated data in the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := testDB.Database("test_database").Collection("questions")
	var dbQuestion map[string]interface{}
	err = collection.FindOne(ctx, map[string]interface{}{"_id": questionID}).Decode(&dbQuestion)
	assert.NoError(t, err)
	assert.Equal(t, updatedTitle, dbQuestion["title"])
	assert.Equal(t, updatedDescription, dbQuestion["description"])
}

func TestDeleteQuestion_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert a test question
	questionID := primitive.NewObjectID()
	testQuestion := utils.GenerateQuestion(map[string]interface{}{"ID": questionID})
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

func TestGetAllQuestions_Integration(t *testing.T) {
	clearDatabase()

	// Arrange: Insert multiple test questions
	questions := []model.Question{
		utils.GenerateQuestion(nil),
		utils.GenerateQuestion(nil),
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
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Len(t, response, len(questions))
}
