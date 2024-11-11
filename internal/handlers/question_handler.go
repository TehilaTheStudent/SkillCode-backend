package handlers

import (
	"fmt"
	"net/http"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
	r.POST("/questions", handler.CreateQuestion)
	r.GET("/questions/:id", handler.GetQuestionByID)
	r.GET("/questions", handler.GetAllQuestions)
	r.PUT("/questions/:id", handler.UpdateQuestion)
	r.DELETE("/questions/:id", handler.DeleteQuestion)
	r.POST("/questions/:id/test", handler.TestQuestion)
}

// CreateQuestion creates a new question
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use go-playground/validator to validate struct fields
	validate := validator.New()
	if err := validate.Struct(question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdQuestion, err := h.Service.CreateQuestion(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdQuestion)
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
	var submission model.Solution
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

func HandleError(c *gin.Context, err error) {
	fmt.Println(err)
	if customErr, ok := err.(*utils.CustomError); ok {
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
