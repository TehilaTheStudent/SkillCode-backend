package service

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
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
	TestQuestion(id string, solution model.Solution) (string, error)
}

type QuestionService struct {
	Repo repository.QuestionRepositoryInterface
}

// NewQuestionService creates a new QuestionService with a QuestionRepository instance.
func NewQuestionService(repo repository.QuestionRepositoryInterface) *QuestionService {
	return &QuestionService{Repo: repo}
}

func (s *QuestionService) CreateQuestion(question model.Question) (*model.Question, error) {
	result, err := s.Repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *QuestionService) GetQuestionByID(id string) (*model.Question, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.New(400, "Invalid ID: "+id)
	}
	return s.Repo.GetQuestionByID(objID)
}

func (s *QuestionService) GetAllQuestions() ([]model.Question, error) {
	return s.Repo.GetAllQuestions()
}

func (s *QuestionService) UpdateQuestion(id string, question model.Question) (*model.Question, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.New(400, "Invalid ID: "+id)
	}
	_, err = s.Repo.UpdateQuestion(objID, question)
	if err != nil {
		return nil, err
	}
	question.ID = objID
	return &question, nil
}

func (s *QuestionService) DeleteQuestion(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.New(400, "Invalid ID: "+id)
	}

	_, err = s.Repo.DeleteQuestion(objID)
	return err
}

// TestQuestion simulates running a user-provided function against test cases for a question.
func (s *QuestionService) TestQuestion(questionId string, solution model.Solution) (string, error) {
	objID, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		return "", utils.New(400, "Invalid ID: "+questionId)
	}
	question, err := s.Repo.GetQuestionByID(objID)
	if err != nil {
		return "", utils.New(404, "Question not found with ID: "+questionId)
	}

	output, err := tester.TestUserSolution(question, solution.Function, solution.Language)
	if err != nil {
		return "", err
	}
	return output, nil
}
