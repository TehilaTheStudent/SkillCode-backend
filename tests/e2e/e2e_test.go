package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/dependencies"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
)

var baseURL string

func TestMain(m *testing.M) {
	dependencies.EnsureWorkingDirectory()
	// Load environment variables
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		if err := godotenv.Load("tests/e2e/.env.e2e"); err != nil {
			log.Fatalf("Error loading .env.e2e file: %v", err)
		}

		// Print loaded environment variables (optional for debugging)
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
	err := json.NewDecoder(retrieveResp.Body).Decode(&retrievedQuestion)
	if err != nil {
		t.Fatalf("Failed to decode retrieved question: %v", err)
	}
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
	err = json.NewDecoder(updatedResp.Body).Decode(&updatedQuestion)
	if err != nil {
		t.Fatalf("Failed to decode updated question: %v", err)
	}
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

func TestUserSolutionSubmission(t *testing.T) {
	// Step 1: Create a question
	createBody := utils.GenerateCreateQuestionPayload(map[string]interface{}{
		"title":      "Binary Search",
		"visibility": "public",
	})
	createResp, createID := createQuestion(t, createBody)
	assert.Equal(t, http.StatusOK, createResp.StatusCode)
	assert.NotEmpty(t, createID)

	// Step 2: Submit a solution for the created question
	solutionBody := `{
		"language": "python",
		"user_function": "def binary_search(arr, target):\n    left, right = 0, len(arr) - 1\n    while left <= right\n        mid = (left + right) // 2\n        if arr[mid] == target\n            return mid\n        elif arr[mid] < target\n            left = mid + 1\n        else\n            right = mid - 1\n    return -1"
	}`

	submitResp := solutionSubmission(t, createID, solutionBody)
	assert.Equal(t, http.StatusOK, submitResp.StatusCode)
	if submitResp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(submitResp.Body)
		assert.NoError(t, err)
		t.Logf("Response body: %s", string(bodyBytes))
	}
	assert.Equal(t, http.StatusOK, submitResp.StatusCode)
	// Step 3: Clean up the created question
	deleteResp := deleteQuestion(t, createID)
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
}

func makeRequest(t *testing.T, method, url, body string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	return resp
}

func createQuestion(t *testing.T, body string) (*http.Response, string) {
	resp := makeRequest(t, "POST", baseURL+"/questions", body)
	defer resp.Body.Close() // Ensure the response body is closed

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	id, _ := response["id"].(string)
	return resp, id
}

func getQuestionByID(t *testing.T, id string) *http.Response {
	return makeRequest(t, "GET", baseURL+"/questions/"+id, "")
}

func updateQuestion(t *testing.T, id string, body string) *http.Response {
	return makeRequest(t, "PUT", baseURL+"/questions/"+id, body)
}

func deleteQuestion(t *testing.T, id string) *http.Response {
	return makeRequest(t, "DELETE", baseURL+"/questions/"+id, "")
}

func solutionSubmission(t *testing.T, questionID string, body string) *http.Response {
	return makeRequest(t, "POST", baseURL+"/questions/"+questionID+"/test", body)
}
