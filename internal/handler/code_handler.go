package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"

	"github.com/gin-gonic/gin"
)

// ServeUtils serves the utilities file based on the language query parameter
func ServeUtils(c *gin.Context) {
	language, err := model.LowerToEnum(c.Query("language"))
	if err != nil {
		LogAndRespondError(c, errors.New("invalid language"), http.StatusBadRequest)
		return
	}
	filePath := config.GlobalConfigSandboxes[language].UtilsFile

	content, err := os.ReadFile(filePath)
	if err != nil {
		LogAndRespondError(c, errors.New("failed to read file: " + filePath), http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "text/plain", content)
}

// RegisterQuestionRoutes sets up the routes for question-related endpoints
func RegisterCodeRoutes(r *gin.Engine) {
	appGroup := r.Group(config.GlobalConfigAPI.Base)
	appGroup.GET("/utils", ServeUtils)
}
