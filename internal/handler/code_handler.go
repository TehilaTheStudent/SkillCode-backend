package handler

import (
	"net/http"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"

	"github.com/gin-gonic/gin"
)

// ServeUtils serves the utilities file based on the language query parameter
func ServeUtils(c *gin.Context) {
	language :=  c.Query("language")
	var filePath string
	
	switch language {
	case "python":
		filePath = config.GlobalConfigSandboxes[model.Python].UtilsFile + ".py"
	case "javascript":
		filePath = config.GlobalConfigSandboxes[model.JavaScript].UtilsFile + ".js"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported language"})
		return
	}

	content, err := os.ReadFile(filePath)
	if (err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	c.Data(http.StatusOK, "text/plain", content)
}

// RegisterQuestionRoutes sets up the routes for question-related endpoints
func RegisterCodeRoutes(r *gin.Engine) {
	appGroup := r.Group(config.GlobalConfigAPI.Base)
	appGroup.GET("/utils", ServeUtils)
}
