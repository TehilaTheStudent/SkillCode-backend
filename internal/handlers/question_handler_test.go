package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockQuestionService struct {
	mock.Mock
}

func (m *MockQuestionService) CreateQuestion(question model.Question) (*model.Question, error) {
	args := m.Called(question)
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockQuestionService) GetQuestionByID(id string) (*model.Question, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockQuestionService) GetAllQuestions() ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *MockQuestionService) UpdateQuestion(id string, question model.Question) (*model.Question, error) {
	args := m.Called(id, question)
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockQuestionService) DeleteQuestion(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionService) TestQuestion(id string, userFunction string) (bool, error) {
	args := m.Called(id, userFunction)
	return args.Bool(0), args.Error(1)
}

func TestCreateQuestion(t *testing.T) {
	// Arrange: Create a mock service and handler
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	// Prepare test data
	title := "Two Sum Problem"
	description := "Find two numbers that add up to a specific target."
	functionSignature := "func twoSum(nums []int, target int) []int"

	// Define the test cases
	testCases := []model.TestCase{
		{
			Input:          []int{2, 7, 11, 15}, // Example input array for the function
			ExpectedOutput: 9,                   // Expected output
		},
		{
			Input:          []int{3, 2, 4}, // Another example input
			ExpectedOutput: 6,              // Expected output
		},
	}

	// Create a new expected question with an ObjectID
	objID := primitive.NewObjectID()
	expectedQuestion := model.Question{
		ID:                objID,
		Title:             title,
		Description:       description,
		FunctionSignature: functionSignature,
		TestCases:         testCases,
		Visibility:        "public",
		CreatedBy:         "user123",
	}

	// Setup mock service to return the created question
	mockService.On("CreateQuestion", mock.AnythingOfType("model.Question")).Return(&expectedQuestion, nil)

	// Create a new Gin engine and register routes
	router := gin.Default()
	RegisterQuestionRoutes(router, handler)

	// Create a mock HTTP request to test the CreateQuestion endpoint
	jsonBody := `{
        "title": "Two Sum Problem",
        "description": "Find two numbers that add up to a specific target.",
        "function_signature": "func twoSum(nums []int, target int) []int",
        "test_cases": [
            {
                "input": [2, 7, 11, 15],
                "expected_output": 9
            },
            {
                "input": [3, 2, 4],
                "expected_output": 6
            }
        ],
        "visibility": "public",
        "created_by": "user123"
    }`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Two Sum Problem")
	assert.Contains(t, w.Body.String(), objID.Hex()) // Assert ID is in string format
	mockService.AssertExpectations(t)
}

func TestUpdateQuestion(t *testing.T) {
	// Arrange
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	// Prepare test data
	objID := primitive.NewObjectID()
	updatedTitle := "Updated Two Sum Problem"
	updatedDescription := "Updated description of the Two Sum Problem"

	// Mock the service
	mockService.On("UpdateQuestion", objID.Hex(), mock.AnythingOfType("model.Question")).
		Return(&model.Question{
			ID:                objID,
			Title:             updatedTitle,
			Description:       updatedDescription,
			FunctionSignature: "func twoSum(nums []int, target int) []int",
			TestCases:         []model.TestCase{},
			Visibility:        "public",
			CreatedBy:         "user123",
		}, nil).Once()

	// Set up router
	router := gin.Default()
	RegisterQuestionRoutes(router, handler)

	// JSON body
	jsonBody := fmt.Sprintf(`{
        "title": "%s",
        "description": "%s",
        "function_signature": "func twoSum(nums []int, target int) []int",
        "test_cases": [],
        "visibility": "public",
        "created_by": "user123"
    }`, updatedTitle, updatedDescription)

	// Execute the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/questions/%s", objID.Hex()), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), updatedTitle)
	assert.Contains(t, w.Body.String(), objID.Hex())
	mockService.AssertExpectations(t)
}

func TestDeleteQuestion(t *testing.T) {
	// Arrange: Create a mock service and handler
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	// Prepare the question ID to delete
	objID := primitive.NewObjectID()

	// Setup mock service to return no error for deleting a question
	mockService.On("DeleteQuestion", objID.Hex()).Return(nil)

	// Create a new Gin engine and register routes
	router := gin.Default()
	RegisterQuestionRoutes(router, handler)

	// Create a mock HTTP request to test the DeleteQuestion endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/questions/"+objID.Hex(), nil)

	// Execute the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Question deleted")
	mockService.AssertExpectations(t)
}
func TestTestQuestion(t *testing.T) {
	// Arrange
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	questionID := primitive.NewObjectID().Hex()
	userFunction := "func twoSum(nums []int, target int) []int { return nil }"
	expectedResult := false

	// Mock the service behavior
	mockService.On("TestQuestion", questionID, userFunction).
		Return(expectedResult, nil).Once()

	// Set up router
	router := gin.Default()
	router.POST("/questions/:id/test", handler.TestQuestion)

	// JSON body
	jsonBody := fmt.Sprintf(`{"user_function": "%s"}`, userFunction)

	// Create mock HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/questions/%s/test", questionID), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, expectedResult, response["result"])
	mockService.AssertExpectations(t)
}

func TestGetAllQuestions(t *testing.T) {
	// Arrange
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	// Test data
	questions := []model.Question{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Two Sum",
			Description: "Find two numbers that add up to a target.",
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

	// Mock the service
	mockService.On("GetAllQuestions").Return(questions, nil).Once()

	// Set up router
	router := gin.Default()
	router.GET("/questions", handler.GetAllQuestions)

	// Create mock HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, len(questions), len(response))
	mockService.AssertExpectations(t)
}

func TestGetQuestionByID(t *testing.T) {
	// Arrange
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	questionID := primitive.NewObjectID()
	question := &model.Question{
		ID:          questionID,
		Title:       "Two Sum",
		Description: "Find two numbers that add up to a target.",
		Visibility:  "public",
		CreatedBy:   "user123",
	}

	// Mock the service
	mockService.On("GetQuestionByID", questionID.Hex()).Return(question, nil).Once()

	// Set up router
	router := gin.Default()
	router.GET("/questions/:id", handler.GetQuestionByID)

	// Create mock HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/questions/%s", questionID.Hex()), nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Question
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, question.ID, response.ID)
	assert.Equal(t, question.Title, response.Title)
	mockService.AssertExpectations(t)
}

func TestGetQuestionByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	questionID := primitive.NewObjectID().Hex()
	fmt.Println(questionID)
	// Mock the service to return "no documents in result"
	mockService.On("GetQuestionByID", questionID).Return(nil, customerrors.New(404, "Question not found with ID:")).Once()

	// Set up router
	router := gin.Default()
	router.GET("/questions/:id", handler.GetQuestionByID)

	// Create mock HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/questions/%s", questionID), nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Question not found")
	mockService.AssertExpectations(t)
}

func TestCreateQuestion_InvalidJSON(t *testing.T) {
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	router := gin.Default()
	router.POST("/questions", handler.CreateQuestion)

	invalidJSON := `title: "Two Sum Problem"`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
