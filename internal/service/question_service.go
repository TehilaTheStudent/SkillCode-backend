package service

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define the service interface
type QuestionServiceInterface interface {
	CreateQuestion(question model.Question) (*model.Question, error)
	GetQuestionByID(id string) (*model.Question, error)
	GetAllQuestions() ([]model.Question, error)
	UpdateQuestion(id string, question model.Question) (*model.Question, error)
	DeleteQuestion(id string) error
	TestQuestion(id string, solution model.Submission) (string, error)
}

type QuestionService struct {
	Repo repository.QuestionRepositoryInterface
}

// NewQuestionService creates a new QuestionService with a QuestionRepository instance.
func NewQuestionService(repo repository.QuestionRepositoryInterface) *QuestionService {
	return &QuestionService{Repo: repo}
}

// CreateQuestion creates a new question in the repository.
func (s *QuestionService) CreateQuestion(question model.Question) (*model.Question, error) {
	result, err := s.Repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handleInvalidID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, utils.New(400, "Invalid ID: "+id)
	}
	return objID, nil
}

// GetQuestionByID retrieves a question by its ID from the repository.
func (s *QuestionService) GetQuestionByID(id string) (*model.Question, error) {
	objID, err := handleInvalidID(id)
	if err != nil {
		return nil, err
	}
	return s.Repo.GetQuestionByID(objID)
}

// GetAllQuestions retrieves all questions from the repository.
func (s *QuestionService) GetAllQuestions() ([]model.Question, error) {
	return s.Repo.GetAllQuestions()
}

// UpdateQuestion updates an existing question in the repository.
func (s *QuestionService) UpdateQuestion(id string, question model.Question) (*model.Question, error) {
	objID, err := handleInvalidID(id)
	if err != nil {
		return nil, err
	}
	_, err = s.Repo.UpdateQuestion(objID, question)
	if err != nil {
		return nil, err
	}
	question.ID = objID
	return &question, nil
}

// DeleteQuestion deletes a question by its ID from the repository.
func (s *QuestionService) DeleteQuestion(id string) error {
	objID, err := handleInvalidID(id)
	if err != nil {
		return err
	}
	_, err = s.Repo.DeleteQuestion(objID)
	return err
}

// TestQuestion simulates running a user-provided function against test cases for a question.
func (s *QuestionService) TestQuestion(questionId string, submission model.Submission) (string, error) {
	objID, err := handleInvalidID(questionId)
	if err != nil {
		return "", err
	}
	//validations:
	question, err := s.Repo.GetQuestionByID(objID)
	if err != nil {
		return "", utils.New(404, "Question not found with ID: "+questionId)
	}
	// Step 3: Prepare Python Test Runner
	sandboxConfig, err := config.NewSandboxConfig(submission.Language)
	if err != nil {
		return "", err
	}
	testRunnerPath := sandboxConfig.UserCodePath
	err = tester.CreatePythonTestRunner(testRunnerPath, *question, submission.Code)
	if err != nil {
		return "", err
	}
	output, err := tester.TestUserSolution(question, submission.Code, submission.Language, *sandboxConfig)
	if err != nil {
		return "", err
	}
	return output, nil
}
