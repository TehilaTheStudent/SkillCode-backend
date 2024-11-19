package handler

import (
	"fmt"
	"net/http"

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

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	// Extract query parameters
	q := c.Query("q")                            // Search query
	categories := c.QueryArray("categories")     // Selected categories
	difficulties := c.QueryArray("difficulties") // Selected difficulties
	sort := c.DefaultQuery("sort", "stats")      // Sorting column
	order := c.DefaultQuery("order", "desc")     // Sorting order

	// Construct query parameters object
	params := model.QuestionQueryParams{
		SearchQuery:  q,
		Categories:   categories,
		Difficulties: difficulties,
		SortField:    sort,
		SortOrder:    order,
	}

	// Call the service layer
	questions, err := h.Service.GetAllQuestions(params)
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
	var submission model.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		LogAndRespondError(c, err, http.StatusBadRequest)
		return
	}
	result, err := h.Service.TestQuestion(id, submission)
	if err != nil {
		LogAndRespondError(c, err, http.StatusInternalServerError)
		return
	}
	// Set Content-Type to application/json and send the result directly
	c.Data(http.StatusOK, "application/json", []byte(result))
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

func LogAndRespondError(c *gin.Context, err error, statusCode int) {
	// Log the error using the middleware
	c.Error(err) // Will be picked up by the middleware

	// Respond with appropriate error message
	if customErr, ok := err.(*model.CustomError); ok {
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
	} else {
		c.JSON(statusCode, gin.H{"error": err.Error()})
	}
}
