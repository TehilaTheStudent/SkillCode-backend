package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/coding"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator"
)

// QuestionHandler holds the service interface
type QuestionHandler struct {
	Service service.QuestionServiceInterface // Use the interface
}

// NewQuestionHandler initializes a QuestionHandler with a given QuestionServiceInterface
func NewQuestionHandler(service service.QuestionServiceInterface) *QuestionHandler {
	return &QuestionHandler{Service: service}
}

// RegisterQuestionRoutes sets up the routes for question-related endpoints
func RegisterQuestionRoutes(r *gin.Engine, handler *QuestionHandler) {
	appGroup := r.Group(config.GlobalConfigAPI.Base)
	appGroup.POST("/questions", handler.CreateQuestion)
	appGroup.GET("/questions/:id", handler.GetQuestionByID)
	appGroup.GET("/questions", handler.GetAllQuestions)
	appGroup.PUT("/questions/:id", handler.UpdateQuestion)
	appGroup.DELETE("/questions/:id", handler.DeleteQuestion)
	appGroup.POST("/questions/:id/test", handler.TestQuestion)
	appGroup.GET("/questions/:id/signature", handler.GetFunctionSignature)
}

// CreateQuestion creates a new question
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		LogAndRespondError(c, err, http.StatusBadRequest)
		return
	}
	// Use go-playground/validator to validate struct fields
	// validate := validator.New()
	// if err := validate.Struct(question); err != nil {
	// 	LogAndRespondError(c, err, http.StatusBadRequest)
	// 	return
	// }
	createdQuestion, err := h.Service.CreateQuestion(question)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, createdQuestion)
}

// GetQuestionByID retrieves a question by its ID
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	id := c.Param("id")
	question, err := h.Service.GetQuestionByID(id)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, question)
}
func splitOrEmpty(value string) []string {
	if value == "" {
		return []string{} // Return an empty slice if the string is empty
	}
	return strings.Split(value, ",")
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	categories := c.Query("categories")
	difficulties := c.Query("difficulties")

	query := model.QuestionQueryParams{
		Search:       c.Query("search"),
		Categories:   splitOrEmpty(categories),
		Difficulties: splitOrEmpty(difficulties),
		SortBy:       c.Query("sort_by"),
		SortOrder:    c.Query("order"),
	}

	// Call the service layer
	questions, err := h.Service.GetAllQuestions(query)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}

	// Respond with the filtered, sorted questions
	c.JSON(http.StatusOK, questions)
}

// UpdateQuestion updates an existing question
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		LogAndRespondError(c, err, http.StatusBadRequest)
		return
	}
	updatedQuestion, err := h.Service.UpdateQuestion(id, question)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, updatedQuestion)
}

// DeleteQuestion deletes a question by its ID
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	err := h.Service.DeleteQuestion(id)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// TestQuestion simulates running a user-provided function against test cases for a question.
func (h *QuestionHandler) TestQuestion(c *gin.Context) {
	id := c.Param("id")

	// Bind the incoming JSON request to the Submission struct
	var submission model.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		LogAndRespondError(c, err, http.StatusBadRequest)
		return
	}

	// Call the service to test the question and get the feedback
	requestID := c.GetString("request_id")
	feedback, err := h.Service.TestUniqueQuestion(id, submission, requestID)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}

	// Serialize the Feedback struct to JSON
	response, err := json.Marshal(feedback)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}

	// Set Content-Type to application/json and send the JSON response
	c.Data(http.StatusOK, "application/json", response)
}

func (h *QuestionHandler) GetFunctionSignature(c *gin.Context) {
	// Extract question ID and language
	id := c.Param("id")
	language := c.Query("language")
	if language == "" {
		LogAndRespondError(c, fmt.Errorf("Language parameter is required"), http.StatusBadRequest)
		return
	}

	// Map lowercase language to PredefinedSupportedLanguage
	langEnum, err := model.LowerToEnum(language)
	if err != nil {
		LogAndRespondError(c, fmt.Errorf("Invalid language parameter"), http.StatusBadRequest)
		return
	}
	// Fetch the question details
	question, err := h.Service.GetQuestionByID(id)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}

	// Generate the function signature
	signature, err := coding.GenerateByQuestionAndLanguage(*question, langEnum)
	if err != nil {
		LogAndRespondError(c, fmt.Errorf("Failed to generate function signature for language %s: %s", language, err.Error()), http.StatusInternalServerError)
		return
	}

	// Return the function signature
	c.JSON(http.StatusOK, gin.H{
		"function_signature": signature,
	})
}
