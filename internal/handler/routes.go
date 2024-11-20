package handler

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// registerRoutes registers all application routes
func RegisterRoutes(r *gin.Engine, questionHandler *QuestionHandler) {

	RegisterQuestionRoutes(r, questionHandler)
	RegisterCodeRoutes(r)
	RegisterConfigRoutes(r)
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
