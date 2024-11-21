package service

import (
	"encoding/json"
	"fmt"
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
	// TestQuestion(id string, solution model.Submission) (*model.Feedback, error)
	TestUniqueQuestion(questionID string, submission model.Submission, requestID string) (*model.Feedback, error)
}

type QuestionService struct {
	Repo         repository.QuestionRepositoryInterface
	SharedTester *tester.SharedTester
}

// NewQuestionService creates a new QuestionService with a QuestionRepository instance.
func NewQuestionService(repo repository.QuestionRepositoryInterface,sharedTester *tester.SharedTester) *QuestionService {
	return &QuestionService{Repo: repo, SharedTester: sharedTester}
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
	if params.Search != "" {
		questions = filterBySearchQuery(questions, params.Search)
	}

	if len(params.Categories) > 0 {
		questions = filterByCategories(questions, params.Categories)
	}

	if len(params.Difficulties) > 0 {
		questions = filterByDifficulties(questions, params.Difficulties)
	}

	// Apply sorting
	questions = sortQuestions(questions, params.SortBy, params.SortOrder)

	// Ensure no nil slices are returned
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
	categorySet := make(map[string]struct{}, len(categories))
	for _, category := range categories {
		if category != "" {
			categorySet[category] = struct{}{}
		}
	}

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
	difficultySet := make(map[string]struct{}, len(difficulties))
	for _, difficulty := range difficulties {
		if difficulty != "" {
			difficultySet[difficulty] = struct{}{}
		}
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
func sortQuestions(questions []model.Question, sortBy, sortOrder string) []model.Question {
	sort.SliceStable(questions, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "stats":
			less = questions[i].Stats < questions[j].Stats
		case "difficulty":
			order := map[string]int{"easy": 1, "medium": 2, "hard": 3}
			less = order[strings.ToLower(questions[i].Difficulty)] < order[strings.ToLower(questions[j].Difficulty)]
		case "category":
			less = strings.ToLower(questions[i].Category) < strings.ToLower(questions[j].Category)
		default: // Default to sorting by title
			less = strings.ToLower(questions[i].Title) < strings.ToLower(questions[j].Title)
		}
		if strings.ToLower(sortOrder) == "desc" {
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

// // TestQuestion simulates running a user-provided function against test cases for a question.
// func (s *QuestionService) TestQuestion(questionId string, submission model.Submission) (*model.Feedback, error) {
// 	objID, err := handleInvalidID(questionId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//validations:
// 	question, err := s.Repo.GetQuestionByID(objID)
// 	if err != nil {
// 		return nil, model.NewCustomError(404, "Question not found with ID: "+questionId)
// 	}
// 	// Increment attempts
// 	question.Stats++
// 	_, err = s.Repo.UpdateQuestion(objID, *question)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sandboxConfig := config.GlobalLanguageConfigs[submission.Language]

// 	testRunnerPath := sandboxConfig.TestUserCodePath
// 	err = tester.CreateTestRunner(submission.Language, testRunnerPath, *question, submission.Code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	output, err := tester.TestUserSolution(question, submission.Code, submission.Language, *sandboxConfig)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if output.Error != nil && *output.Error == model.InternalServerError {
// 		if output.Details == nil {
// 			return nil, model.NewCustomError(500, "Internal Server Error")
// 		}
// 		return nil, model.NewCustomError(500, "Internal Server Error: "+*output.Details)
// 	}

// 	return output, nil
// }

func (s *QuestionService) TestUniqueQuestion(questionID string, submission model.Submission, requestID string) (*model.Feedback, error) {
	// Step 1: Validate Question
	objID, err := handleInvalidID(questionID)
	if err != nil {
		return nil, err
	}

	question, err := s.Repo.GetQuestionByID(objID)
	if err != nil {
		return nil, model.NewCustomError(404, "Question not found with ID: "+questionID)
	}

	// Step 2: Create unique assets
	uniqueDir, err := tester.GenerateUniqueAssets(requestID, *question, submission)
	if err != nil {
		return nil, err
	}
	defer tester.CleanupUniqueAssets(uniqueDir) // Ensure cleanup

	// Step 3: Create UniqueTester and execute
	uniqueTester := tester.NewUniqueTester(
		s.SharedTester,
		fmt.Sprintf("job-%s", requestID),
		config.GlobalLanguageConfigs[submission.Language].ImageName,
		model.GetRuntime(submission.Language),
		model.GetFileExtension(submission.Language),
		requestID,
	)

	rawLogs, err := uniqueTester.ExecuteUniqueTest(uniqueDir, submission.Code)
	if err != nil {
		return nil, err
	}
	// Parse JSON logs into Feedback struct
	var feedback model.Feedback
	if parseErr := json.Unmarshal([]byte(rawLogs), &feedback); parseErr != nil {
		return nil, model.NewCustomError(500, fmt.Sprintf("failed to parse feedback logs: %v", parseErr))
	}

	// Step 4: Process results
	if feedback.Error != nil && *feedback.Error == model.InternalServerError {
		if feedback.Details == nil {
			return nil, model.NewCustomError(500, "Internal Server Error")
		}
		return nil, model.NewCustomError(500, "Internal Server Error: "+*feedback.Details)
	}

	return &feedback, nil
}
