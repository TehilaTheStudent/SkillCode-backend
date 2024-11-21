package model

import "encoding/json"

type ErrorType string

const (
    CompilationError       ErrorType = "compilation"
    FailTestsError         ErrorType = "fail tests"
    InternalServerError    ErrorType = "internal server error"
)

type Feedback struct {
    Status  string    `json:"status"`            // Overall status: success or fail
    Results []Result  `json:"results"`           // Array of individual test case results
    Error   *ErrorType `json:"error,omitempty"`   // Error type: compilation, fail tests, internal server error, or null
    Details *string   `json:"details,omitempty"` // Detailed error description, or null if not applicable
}

type Result struct {
    Status        string          `json:"status"`         // Status of the test case: pass or fail
    Parameters    []string        `json:"parameters"`     // Array of strings representing input parameters
    ExpectedOutput json.RawMessage `json:"expected_output"` // Expected output (can be a string or number)
    ActualOutput   json.RawMessage `json:"actual_output"`   // Actual output (can be a string or number)
}

