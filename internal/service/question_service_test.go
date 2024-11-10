package service

import (
	"testing"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repository for unit testing
type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) CreateQuestion(question model.Question) (*model.Question, error) {
	args := m.Called(question)
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockQuestionRepository) GetQuestionByID(id primitive.ObjectID) (*model.Question, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockQuestionRepository) GetAllQuestions() ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *MockQuestionRepository) UpdateQuestion(id primitive.ObjectID, question model.Question) (bool, error) {
	args := m.Called(id, question)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuestionRepository) DeleteQuestion(id primitive.ObjectID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestCreateQuestion(t *testing.T) {
	// Arrange: Create a mock repository and service
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Prepare a question object with optional fields as pointers
	title := "Test Question"
	description := "Test description"
	functionSignature := "func test()"
	question := model.Question{
		Title:             title,
		Description:       description,
		FunctionSignature: functionSignature,
	}

	// Create a new expected question with an ObjectID
	objID := primitive.NewObjectID()
	expectedQuestion := model.Question{
		ID:                objID,
		Title:             title,
		Description:       description,
		FunctionSignature: functionSignature,
	}

	// Setup the mock behavior
	mockRepo.On("CreateQuestion", question).Return(&expectedQuestion, nil)

	// Act: Call the service method
	result, err := service.CreateQuestion(question)

	// Assert: Verify the result
	assert.NoError(t, err)
	assert.Equal(t, expectedQuestion.ID, result.ID)
	assert.Equal(t, expectedQuestion.Title, result.Title)
	assert.Equal(t, expectedQuestion.Description, result.Description)
	assert.Equal(t, expectedQuestion.FunctionSignature, result.FunctionSignature)

	// Verify that the mock repository was called as expected
	mockRepo.AssertExpectations(t)
}

func TestGetQuestionByID(t *testing.T) {
	// Arrange: Create a mock repository and service
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Create a question with ID
	questionID := "672fef1c2874f7d1ff6cff66"
	objID, _ := primitive.ObjectIDFromHex(questionID)
	expectedQuestion := model.Question{
		ID:                objID,
		Title:             "Test Question",
		Description:       "nil",
		FunctionSignature: "nil",
	}

	// Setup the mock behavior
	mockRepo.On("GetQuestionByID", objID).Return(&expectedQuestion, nil)

	// Act: Call the service method
	result, err := service.GetQuestionByID(questionID)

	// Assert: Verify the result
	assert.NoError(t, err)
	assert.Equal(t, expectedQuestion.ID, result.ID)
	assert.Equal(t, expectedQuestion.Title, result.Title)

	// Verify that the mock repository was called as expected
	mockRepo.AssertExpectations(t)
}

func TestUpdateQuestion(t *testing.T) {
	// Arrange: Create a mock repository and service
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Prepare a question object for updating
	title := "title: updated"
	updatedQuestion := model.Question{
		ID:                primitive.NewObjectID(),
		Title:             title,
		Description:       "nil",
		FunctionSignature: "nil",
	}

	// Mock the repository to return success for the update
	mockRepo.On("UpdateQuestion", updatedQuestion.ID, updatedQuestion).Return(true, nil)

	// Act: Call the service method
	result, err := service.UpdateQuestion(updatedQuestion.ID.Hex(), updatedQuestion)

	// Assert: Verify the result
	assert.NoError(t, err)
	assert.Equal(t, updatedQuestion.ID, result.ID)
	assert.Equal(t, updatedQuestion.Title, result.Title)

	// Verify that the mock repository was called as expected
	mockRepo.AssertExpectations(t)
}

func TestDeleteQuestion(t *testing.T) {
	// Arrange: Create a mock repository and service
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Setup the mock behavior for DeleteQuestion
	mockRepo.On("DeleteQuestion", mock.Anything).Return(true, nil)

	// Act: Call the service method
	err := service.DeleteQuestion("672fef1c2874f7d1ff6cff66")

	// Assert: Verify the result
	assert.NoError(t, err)

	// Verify that the mock repository was called as expected
	mockRepo.AssertExpectations(t)
}

func TestGetQuestionByID_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Act: Call the service method with an invalid ID
	_, err := service.GetQuestionByID("invalid-id")

	// Assert: Verify the result
	assert.Error(t, err)
	assert.IsType(t, &customerrors.CustomError{}, err)
	customErr := err.(*customerrors.CustomError)
	assert.Equal(t, 400, customErr.Code)
	assert.Equal(t, "Invalid ID: invalid-id", customErr.Message)
}
func TestUpdateQuestion_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	updatedQuestion := model.Question{
		Title:             "Updated Test Question",
		Description:       "This is an updated test question",
		FunctionSignature: "func updatedTest() {}",
		TestCases: []model.TestCase{
			{Input: "input1", ExpectedOutput: "output1"},
		},
		Visibility: "public",
		CreatedBy:  "tester",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act: Call the service method with an invalid ID
	_, err := service.UpdateQuestion("invalid-id", updatedQuestion)

	// Assert: Verify the result
	assert.Error(t, err)
	assert.IsType(t, &customerrors.CustomError{}, err)
	customErr := err.(*customerrors.CustomError)
	assert.Equal(t, 400, customErr.Code)
	assert.Equal(t, "Invalid ID: invalid-id", customErr.Message)
}

func TestDeleteQuestion_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Act: Call the service method with an invalid ID
	err := service.DeleteQuestion("invalid-id")

	// Assert: Verify the result
	assert.Error(t, err)
	assert.IsType(t, &customerrors.CustomError{}, err)
	customErr := err.(*customerrors.CustomError)
	assert.Equal(t, 400, customErr.Code)
	assert.Equal(t, "Invalid ID: invalid-id", customErr.Message)
}
