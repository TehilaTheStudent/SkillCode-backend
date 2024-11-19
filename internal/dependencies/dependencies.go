package dependencies

// SetupAllDependencies sets up all necessary dependencies for the application
func SetupAllDependencies() error {
	// Ensure the working directory is correct
	err := EnsureWorkingDirectory()
	if err != nil {
		return err
	}

	// Setup submission dependencies
	err = SetupSubmissionDependencies()
	if err != nil {
		return err
	}

	// Initialize the database connection and start health checks
	err = InitializeDatabase()
	if err != nil {
		return err
	}

	return nil
}

