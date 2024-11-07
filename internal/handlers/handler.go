package handlers

import (
	"net/http"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/models"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterQuestionRoutes(r *gin.Engine) {
	r.POST("/questions", CreateQuestion)
	r.GET("/questions/:id", GetQuestionByID)
	r.GET("/questions", GetAllQuestions)
	r.PUT("/questions/:id", UpdateQuestion)
	r.DELETE("/questions/:id", DeleteQuestion)
	r.POST("/questions/:id/test", TestQuestion)
}

func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdQuestion, err := services.CreateQuestion(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdQuestion)
}

func GetQuestionByID(c *gin.Context) {
	id := c.Param("id")
	question, err := services.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}

func GetAllQuestions(c *gin.Context) {
	questions, err := services.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedQuestion, err := services.UpdateQuestion(id, question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedQuestion)
}

func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

func TestQuestion(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request struct {
		UserFunction string `json:"user_function"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.TestQuestion(objID, request.UserFunction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}
