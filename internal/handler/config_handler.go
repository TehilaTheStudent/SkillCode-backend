package handler

import (
	"net/http"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// ServeConfigs serves the enums and predefined configurations for the frontend
func ServeConfigs(c *gin.Context) {
	config := map[string]interface{}{
		"atomicTypes": []string{
			string(model.Integer),
			string(model.Double),
			string(model.String),
			string(model.Boolean),
		},
		"compositeTypes": []string{
			string(model.GraphNode),
			string(model.TreeNode),
			string(model.ListNode),
			string(model.Array),
			string(model.Matrix),
		},
		"difficulties": []string{
			string(model.Easy),
			string(model.Medium),
			string(model.Hard),
		},
		"categories": []string{
			string(model.ArrayCategory),
			string(model.GraphCategory),
			string(model.StringCategory),
			string(model.TreeCategory),
			string(model.DynamicProgrammingCategory),
			string(model.LinkedListCategory),
			string(model.MatrixCategory),
		},
		"languages": []string{
			string(model.JavaScript),
			string(model.Python),
			string(model.Java),
			string(model.Go),
			string(model.CSharp),
			string(model.Cpp),
		},
	}

	c.JSON(http.StatusOK, config)
}

func RegisterConfigRoutes(r *gin.Engine) {
	r.GET("/configs", ServeConfigs)
}
