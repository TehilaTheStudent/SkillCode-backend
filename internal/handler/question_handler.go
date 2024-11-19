package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/coding"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	appGroup := r.Group(config.LoadConfigAPI().Base)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use go-playground/validator to validate struct fields
	// validate := validator.New()
	// if err := validate.Struct(question); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	createdQuestion, err := h.Service.CreateQuestion(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdQuestion)
}

// GetQuestionByID retrieves a question by its ID
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	id := c.Param("id")
	question, err := h.Service.GetQuestionByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	// Extract query parameters
	q := c.Query("q")                               // Search query
	categories := c.QueryArray("categories")        // Selected categories
	difficulties := c.QueryArray("difficulties")    // Selected difficulties
	sort := c.DefaultQuery("sort", "stats")         // Sorting column
	order := c.DefaultQuery("order", "desc")        // Sorting order

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedQuestion, err := h.Service.UpdateQuestion(id, question)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedQuestion)
}

// DeleteQuestion deletes a question by its ID
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	err := h.Service.DeleteQuestion(id)

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// TestQuestion simulates running a user-provided function against test cases for a question.
func (h *QuestionHandler) TestQuestion(c *gin.Context) {
	id := c.Param("id")
	var submission model.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.Service.TestQuestion(id, submission)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": result})
}

func (h *QuestionHandler) GetFunctionSignature(c *gin.Context) {
	// Extract question ID and language
	id := c.Param("id")
	language := c.Query("language")
	if language == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language parameter is required"})
		return
	}

	// Map lowercase language to PredefinedSupportedLanguage
	langEnum, err := utils.LowerToEnum(language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language parameter"})
		return
	}
	// Fetch the question details
	question, err := h.Service.GetQuestionByID(id)
	if err != nil {
		log.Printf("Error fetching question ID %s: %v", id, err)
		HandleError(c, err)
		return
	}

	// Generate the function signature
	signature, err := coding.GenerateByQuestionAndLanguage(*question, langEnum)
	if err != nil {
		log.Printf("Error generating function signature for question ID %s and language %s: %v", id, language, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to generate function signature for language %s: %s", language, err.Error()),
		})
		return
	}

	// Return the function signature
	c.JSON(http.StatusOK, gin.H{
		"function_signature": signature,
	})
}

func HandleError(c *gin.Context, err error) {
	// Retrieve the logger from the context
	logger, exists := c.Get("logger")
	if !exists {
		// Fallback: Use a default logger
		logger = zap.L()
	}

	// Log the error with contextual information
	LogError(c, logger.(*zap.Logger), err, http.StatusInternalServerError, err.Error())

	// Respond with appropriate error message
	if customErr, ok := err.(*utils.CustomError); ok {
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

func LogError(c *gin.Context, logger *zap.Logger, err error, statusCode int, message string) {
	logger.Error("Error occurred",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", statusCode),
		zap.String("client_ip", c.ClientIP()),
		zap.Error(err),
	)
}
