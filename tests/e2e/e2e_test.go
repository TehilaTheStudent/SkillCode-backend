package e2e

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
)

var baseURL string

func TestMain(m *testing.M) {
	utils.EnsureWorkingDirectory()
	// Load environment variables
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		if err := godotenv.Load("tests/e2e/.env.e2e"); err != nil {
			log.Fatalf("Error loading .env.e2e file: %v", err)
		}

		// Print loaded environment variables (optional for debugging)
		log.Printf("Loaded environment: APP_URL=%s", os.Getenv("APP_URL"))
		baseURL = os.Getenv("APP_URL")
	}

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

func TestCreateRetrieveUpdateDeleteQuestion(t *testing.T) {
	// Step 1: Create a new question
	createBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":      "Merge Sort",
		"visibility": "public",
	})
	createResp, createID := createQuestion(t, createBody)
	assert.Equal(t, http.StatusOK, createResp.StatusCode)
	assert.NotEmpty(t, createID)

	// Step 2: Retrieve the created question
	retrieveResp := getQuestionByID(t, createID)
	assert.Equal(t, http.StatusOK, retrieveResp.StatusCode)

	var retrievedQuestion map[string]interface{}
	json.NewDecoder(retrieveResp.Body).Decode(&retrievedQuestion)
	assert.Equal(t, "Merge Sort", retrievedQuestion["title"])

	// Step 3: Update the question
	updateBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":      "Updated Merge Sort",
		"visibility": "private",
	})
	updateResp := updateQuestion(t, createID, updateBody)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	// Verify the update
	updatedResp := getQuestionByID(t, createID)
	assert.Equal(t, http.StatusOK, updatedResp.StatusCode)

	var updatedQuestion map[string]interface{}
	json.NewDecoder(updatedResp.Body).Decode(&updatedQuestion)
	assert.Equal(t, "Updated Merge Sort", updatedQuestion["title"])
	assert.Equal(t, "private", updatedQuestion["visibility"])

	// Step 4: Delete the question
	deleteResp := deleteQuestion(t, createID)
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)

	// Verify deletion
	notFoundResp := getQuestionByID(t, createID)
	assert.Equal(t, http.StatusNotFound, notFoundResp.StatusCode)
}

func TestErrorHandlingOnInvalidInput(t *testing.T) {
	// Step 1: Test creating a question with missing required fields
	invalidCreateBody := "{\"title\": \"Invalid Question\"}"
	createResp, _ := createQuestion(t, invalidCreateBody)
	assert.Equal(t, http.StatusBadRequest, createResp.StatusCode)

	// Step 2: Test retrieving a non-existent question
	nonExistentID := "000000000000000000000000"
	retrieveResp := getQuestionByID(t, nonExistentID)
	assert.Equal(t, http.StatusNotFound, retrieveResp.StatusCode)

	// Step 3: Test updating a non-existent question
	updateBody := utils.GenerateCreateQuestionPayload(nil)
	updateResp := updateQuestion(t, nonExistentID, updateBody)
	assert.Equal(t, http.StatusNotFound, updateResp.StatusCode)

	// Step 4: Test deleting a non-existent question
	deleteResp := deleteQuestion(t, nonExistentID)
	assert.Equal(t, http.StatusNotFound, deleteResp.StatusCode)
}

func createQuestion(t *testing.T, body string) (*http.Response, string) {
	client := &http.Client{}

	// Use body directly as it's already a JSON string
	req, err := http.NewRequest("POST", baseURL+"/questions", bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close() // Ensure the response body is closed

	// Decode the response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	// Extract the ID safely
	id, _ := response["id"].(string)

	return resp, id
}

func getQuestionByID(t *testing.T, id string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"/questions/"+id, nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

func updateQuestion(t *testing.T, id string, body string) *http.Response {
	client := &http.Client{}

	// Use the JSON string directly as the request body
	req, err := http.NewRequest("PUT", baseURL+"/questions/"+id, bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

func deleteQuestion(t *testing.T, id string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", baseURL+"/questions/"+id, nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}
