package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
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

func (m *MockQuestionService) TestQuestion(id string, solution model.Solution) (string, error) {
	args := m.Called(id, solution)
	return args.String(0), args.Error(1)
}

func TestCreateQuestion(t *testing.T) {
	// Arrange: Create a mock service and handler
	mockService := new(MockQuestionService)
	handler := NewQuestionHandler(mockService)

	// Create a new expected question with an ObjectID and timestamps
	objID := primitive.NewObjectID()
	expectedQuestion := utils.GenerateQuestion(map[string]interface{}{
		"ID":          objID,
		"Title":       "Two Sum Problem",
		"Description": "Find two numbers that add up to a specific target.",
	})

	// Setup mock service to return the created question
	mockService.On("CreateQuestion", mock.AnythingOfType("model.Question")).Return(&expectedQuestion, nil)

	// Create a new Gin engine and register routes
	router := gin.Default()
	RegisterQuestionRoutes(router, handler)

	// Generate JSON payload
	jsonBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":       "Two Sum Problem",
		"description": "Find two numbers that add up to a specific target.",
		"test_cases": []map[string]string{
			{"input": "[2, 7, 11, 15]", "expected_output": "9"},
			{"input": "[3, 2, 4]", "expected_output": "6"},
		},
		"languages": []map[string]string{
			{"language": "golang", "function_signature": "func twoSum(nums []int, target int) []int"},
			{"language": "python", "function_signature": "def two_sum(nums: List[int], target: int) -> List[int]:"},
		},
		"visibility": "public",
		"created_by": "user123",
	})

	// Act: Create a mock HTTP request to test the CreateQuestion endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert: Validate the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Two Sum Problem") // Check title
	assert.Contains(t, w.Body.String(), objID.Hex())       // Check ID is returned
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
			ID:          objID,
			Title:       updatedTitle,
			Description: updatedDescription,
			TestCases:   []model.InputOutput{},
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

		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Reverse String",
			Description: "Reverse a given string.",

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
	mockService.On("GetQuestionByID", questionID).Return(nil, utils.New(404, "Question not found with ID:")).Once()

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
