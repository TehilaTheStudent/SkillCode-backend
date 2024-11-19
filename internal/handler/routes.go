package handler

import (
	"github.com/gin-gonic/gin"
)


// registerRoutes registers all application routes
func RegisterRoutes(r *gin.Engine, questionHandler *QuestionHandler) {

	RegisterQuestionRoutes(r, questionHandler)
	RegisterCodeRoutes(r)
	RegisterConfigRoutes(r)
}
