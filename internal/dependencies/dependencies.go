package dependencies

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/tester"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAllDependencies sets up all necessary dependencies for the application
func SetupAllDependencies() (*mongo.Client, *tester.SharedTester, error) {

	// Setup submission dependencies
	sharedTester, err := SetupSubmissionDependencies(config.GlobalConfigAPI.KubeconfigPath, config.GlobalConfigAPI.Namespace)
	if err != nil {
		return nil, nil, err
	}

	// Initialize the database connection and start health checks
	client, err := InitializeDatabase()
	if err != nil {
		return nil, nil, err
	}

	return client, sharedTester, nil
}
