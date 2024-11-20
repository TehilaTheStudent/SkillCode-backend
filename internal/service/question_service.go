package service

import (
	"sort"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define the service interface
type QuestionServiceInterface interface {
	CreateQuestion(question model.Question) (*model.Question, error)
	GetQuestionByID(id string) (*model.Question, error)
	GetAllQuestions(params model.QuestionQueryParams) ([]model.Question, error)
	UpdateQuestion(id string, question model.Question) (*model.Question, error)
	DeleteQuestion(id string) error
	TestQuestion(id string, solution model.Submission) (*model.Feedback, error)
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
		return primitive.NilObjectID, model.NewCustomError(400, "Invalid ID: "+id)
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

func (s *QuestionService) GetAllQuestions(params model.QuestionQueryParams) ([]model.Question, error) {
	// Fetch all questions from the repository
	questions, err := s.Repo.GetAllQuestions()
	if err != nil {
		return []model.Question{}, err
	}

	// Apply filtering
	if params.SearchQuery != "" {
		questions = filterBySearchQuery(questions, params.SearchQuery)
	}

	if len(params.Categories) > 0 {
		questions = filterByCategories(questions, params.Categories)
	}

	if len(params.Difficulties) > 0 {
		questions = filterByDifficulties(questions, params.Difficulties)
	}

	// Apply sorting
	questions = sortQuestions(questions, params.SortField, params.SortOrder)
	if questions == nil {
		questions = []model.Question{}
	}
	return questions, nil
}

// Filter by search query
func filterBySearchQuery(questions []model.Question, query string) []model.Question {
	var filtered []model.Question
	query = strings.ToLower(query)
	for _, question := range questions {
		if strings.Contains(strings.ToLower(question.Title), query) {
			filtered = append(filtered, question)
		}
	}
	return filtered
}

// Filter by categories
func filterByCategories(questions []model.Question, categories []string) []model.Question {
	categorySet := make(map[string]struct{})
	for _, category := range categories {
		if category != "" {
			categorySet[category] = struct{}{}
		}
	}

	if len(categorySet) == 0 {
		return questions
	}
	return questions

	var filtered []model.Question
	for _, question := range questions {
		if _, exists := categorySet[question.Category]; exists {
			filtered = append(filtered, question)
		}
	}
	return filtered
}

// Filter by difficulties
func filterByDifficulties(questions []model.Question, difficulties []string) []model.Question {
	difficultySet := make(map[string]struct{})
	for _, difficulty := range difficulties {
		if difficulty != "" {
			difficultySet[difficulty] = struct{}{}
		}
	}
	if len(difficultySet) == 0 {
		return questions
	}
	var filtered []model.Question
	for _, question := range questions {
		if _, exists := difficultySet[question.Difficulty]; exists {
			filtered = append(filtered, question)
		}
	}
	return filtered
}

// Sort questions
func sortQuestions(questions []model.Question, sortField, sortOrder string) []model.Question {
	sort.Slice(questions, func(i, j int) bool {
		var less bool
		switch sortField {
		case "stats":
			less = questions[i].Stats < questions[j].Stats
		case "difficulty":
			order := map[string]int{"Easy": 1, "Medium": 2, "Hard": 3}
			less = order[questions[i].Difficulty] < order[questions[j].Difficulty]
		case "category":
			less = questions[i].Category < questions[j].Category
		default:
			less = questions[i].Title < questions[j].Title
		}
		if sortOrder == "desc" {
			return !less
		}
		return less
	})
	return questions
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
func (s *QuestionService) TestQuestion(questionId string, submission model.Submission) (*model.Feedback, error) {
	objID, err := handleInvalidID(questionId)
	if err != nil {
		return nil, err
	}
	//validations:
	question, err := s.Repo.GetQuestionByID(objID)
	if err != nil {
		return nil, model.NewCustomError(404, "Question not found with ID: "+questionId)
	}
	// Step 3: Prepare Python Test Runner
	sandboxConfig := config.GlobalConfigSandboxes[submission.Language]

	testRunnerPath := sandboxConfig.TestUserCodePath
	err = tester.CreateTestRunner(submission.Language, testRunnerPath, *question, submission.Code)
	if err != nil {
		return nil, err
	}
	output, err := tester.TestUserSolution(question, submission.Code, submission.Language, *sandboxConfig)
	if err != nil {
		return nil, err
	}
	return output, nil
}
