package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/parser_validator"

	"github.com/gin-gonic/gin"
)

// ServeUtils serves the utilities file based on the language query parameter
func ServeUtils(c *gin.Context) {
	language, err := model.LowerToEnum(c.Query("language"))
	if err != nil {

		LogAndRespondError(c, errors.New("invalid language"), http.StatusBadRequest)
		return
	}
	filePath := config.GlobalLanguageConfigs[language].UtilsFile

	content, err := os.ReadFile(filePath)
	if err != nil {
		LogAndRespondError(c, errors.New("failed to read file: "+filePath), http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "text/plain", content)
}

// Handler to generate an example string for AbstractType
func GenerateExampleHandler(c *gin.Context) {
	var abstractType model.AbstractType

	// Parse JSON input into AbstractType
	if err := c.ShouldBindJSON(&abstractType); err != nil {
		LogAndRespondError(c, errors.New("invalid input"+err.Error()), http.StatusBadRequest)
		return
	}

	// Generate a valid string example
	example := parser_validator.GenerateValidString(&abstractType)

	// Return the generated string
	c.JSON(http.StatusOK, gin.H{"example": example})
}

// RegisterQuestionRoutes sets up the routes for question-related endpoints
func RegisterCodeRoutes(r *gin.Engine) {
	appGroup := r.Group(config.GlobalConfigAPI.Base)
	appGroup.GET("/ds_utils", ServeUtils)
	appGroup.POST("/ds_utils/examples", GenerateExampleHandler)
}
