package service

import (
	"testing"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
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
	// Create a new expected question with an ObjectID
	objID := primitive.NewObjectID()
	question := utils.GenerateQuestion(map[string]interface{}{"ID": objID})

	expectedQuestion := utils.GenerateQuestion(map[string]interface{}{"ID": objID})

	// Setup the mock behavior
	mockRepo.On("CreateQuestion", question).Return(&question, nil)

	// Act: Call the service method
	result, err := service.CreateQuestion(question)

	// Assert: Verify the result
	assert.NoError(t, err)
	assert.Equal(t, expectedQuestion.ID, result.ID)
	assert.Equal(t, expectedQuestion.Title, result.Title)
	assert.Equal(t, expectedQuestion.Description, result.Description)

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
	expectedQuestion := utils.GenerateQuestion(map[string]interface{}{"ID": objID})

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
	updatedQuestion := utils.GenerateQuestion(map[string]interface{}{"title": title})

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

func assertInvalidIDError(t *testing.T, err error, invalidID string) {
	assert.Error(t, err)
	assert.IsType(t, &utils.CustomError{}, err)
	customErr := err.(*utils.CustomError)
	assert.Equal(t, 400, customErr.Code)
	assert.Equal(t, "Invalid ID: "+invalidID, customErr.Message)
}

func TestGetQuestionByID_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Act: Call the service method with an invalid ID
	_, err := service.GetQuestionByID("invalid-id")

	// Assert: Verify the result
	assertInvalidIDError(t, err, "invalid-id")
}

func TestUpdateQuestion_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	updatedQuestion := model.Question{
		Title:       "Updated Test Question",
		Description: "This is an updated test question",
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
	assertInvalidIDError(t, err, "invalid-id")
}

func TestDeleteQuestion_invalidId(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)

	// Act: Call the service method with an invalid ID
	err := service.DeleteQuestion("invalid-id")

	// Assert: Verify the result
	assertInvalidIDError(t, err, "invalid-id")
}
