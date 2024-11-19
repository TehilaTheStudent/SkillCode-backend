package handler

import (
	"net/http"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"

	"github.com/gin-gonic/gin"
)

// ServePythonUtils serves the Python utilities file
func ServePythonUtils(c *gin.Context) {
	serveFileContent(c, config.NewConfigCode(model.Python).UtilsFile+".py")
}

// ServeJSUtils serves the JavaScript utilities file
func ServeJSUtils(c *gin.Context) {
	serveFileContent(c, config.NewConfigCode(model.JavaScript).UtilsFile+".js")
}

// ServeUtils serves the utilities file based on the language query parameter
func ServeUtils(c *gin.Context) {
	language := c.Query("language")
	var filePath string

	switch language {
	case "python":
		filePath = config.NewConfigCode(model.Python).UtilsFile + ".py"
	case "javascript":
		filePath = config.NewConfigCode(model.JavaScript).UtilsFile + ".js"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported language"})
		return
	}

	serveFileContent(c, filePath)
}

func serveFileContent(c *gin.Context, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	c.Data(http.StatusOK, "text/plain", content)
}

// RegisterQuestionRoutes sets up the routes for question-related endpoints
func RegisterCodeRoutes(r *gin.Engine) {
	appGroup := r.Group(config.LoadConfigAPI().Base)
	appGroup.GET("/utils", ServeUtils)

}
