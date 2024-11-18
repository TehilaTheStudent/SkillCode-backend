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

// GetAllQuestions retrieves all questions
func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	questions, err := h.Service.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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

	// Validate the language
	validLanguages := map[string]bool{
		"Python":     true,
		"JavaScript": true,
		"Java":       true,
		"Go":         true,
		"CSharp":     true,
		"Cpp":        true,
	}
	if !validLanguages[language] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported language: %s", language)})
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
	// Convert language string to PredefinedSupportedLanguage type
	langEnum := model.PredefinedSupportedLanguage(language)
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
	fmt.Println(err)
	if customErr, ok := err.(*utils.CustomError); ok {
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
