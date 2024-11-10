package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
    "log"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
)

var baseURL string

func TestMain(m *testing.M) {
	
	// Load environment variables
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		if err := godotenv.Load(".env.e2e"); err != nil {
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
	createBody := map[string]interface{}{
		"title":              "Merge Sort",
		"description":        "Implement merge sort algorithm",
		"function_signature": "func mergeSort(nums []int) []int",
		"test_cases":         []interface{}{},
		"visibility":         "public",
		"created_by":         "user123",
	}
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
	updateBody := map[string]interface{}{
		"title":              "Updated Merge Sort",
		"description":        "Updated description for merge sort",
		"function_signature": "func mergeSort(nums []int) []int",
		"visibility":         "private",
	}
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
	invalidCreateBody := map[string]interface{}{
		"title": "Invalid Question",
	}
	createResp, _ := createQuestion(t, invalidCreateBody)
	assert.Equal(t, http.StatusBadRequest, createResp.StatusCode)

	// Step 2: Test retrieving a non-existent question
	nonExistentID := "000000000000000000000000"
	retrieveResp := getQuestionByID(t, nonExistentID)
	assert.Equal(t, http.StatusNotFound, retrieveResp.StatusCode)

	// Step 3: Test updating a non-existent question
	updateBody := map[string]interface{}{
		"title": "Should Not Exist",
	}
	updateResp := updateQuestion(t, nonExistentID, updateBody)
	assert.Equal(t, http.StatusNotFound, updateResp.StatusCode)

	// Step 4: Test deleting a non-existent question
	deleteResp := deleteQuestion(t, nonExistentID)
	assert.Equal(t, http.StatusNotFound, deleteResp.StatusCode)
}

func createQuestion(t *testing.T, body map[string]interface{}) (*http.Response, string) {
	client := &http.Client{}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", baseURL+"/questions", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	id := ""
	if resp.StatusCode == http.StatusOK {
		id = response["id"].(string)
	}

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

func updateQuestion(t *testing.T, id string, body map[string]interface{}) *http.Response {
	client := &http.Client{}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("PUT", baseURL+"/questions/"+id, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

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
