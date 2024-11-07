package main

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handlers"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default MongoDB URI
	}
	db.ConnectMongoDB(mongoURI)

	// Initialize the question collection
	repository.InitQuestionCollection()

	// Gin setup
	r := gin.Default()
	handlers.RegisterQuestionRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
