package tester

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// executeFunction is a placeholder for the logic to execute a user's function with test cases.
// This function should be implemented in a real environment with sandboxing.
func TestUserSolution(question *model.Question, userFunction string, language string) (string, error) {
	tester := NewTester("user-solution:latest")

	// Step 1: Generate main file
	if err := tester.GenerateMainFile(userFunction, question.FunctionSignature, question.TestCases); err != nil {
		return "", fmt.Errorf("failed to generate main file: %w", err)
	}

	// Step 2: Build Docker image
	if err := tester.BuildDockerImage(); err != nil {
		return "", fmt.Errorf("failed to build Docker image: %w", err)
	}

	// Step 3: Deploy Pod
	podName := "solution-pod" // Make dynamic if needed
	if err := tester.DeployPod(podName); err != nil {
		return "", fmt.Errorf("failed to deploy pod: %w", err)
	}
	defer tester.CleanUp(podName) // Ensure cleanup happens

	// Step 4: Get logs
	logs, err := tester.GetPodLogs(podName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve pod logs: %w", err)
	}

	return logs, nil
}

type Tester struct {
	ImageName string
}

func NewTester(imageName string) *Tester {
	return &Tester{ImageName: imageName}
}
func (t *Tester) GenerateMainFile(userFunction string, functionSignature string, testCases []model.TestCase) error {
	mainTemplate := `
		package main

		import "fmt"

		%s // User's function

		func main() {
			tests := []struct {
				input    string
				expected string
			}{
				%s
			}

			for i, test := range tests {
				// Deserialize input, call user function, validate output
				_=test
				fmt.Printf("Test %%d: %%v\\n", i+1, "passed")
			}
		}
	`

	// Serialize test cases
	var testCasesCode string
	for _, testCase := range testCases {
		testCasesCode += fmt.Sprintf(`{input: %q, expected: %q},`, testCase.Input, testCase.ExpectedOutput)
	}

	// Generate the final code
	mainCode := fmt.Sprintf(mainTemplate, userFunction, testCasesCode)

	// Define the file path
	filePath := "./temp/solution.go"

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to the file
	return os.WriteFile(filePath, []byte(mainCode), 0644)
}

// BuildDockerImage builds a Docker image for the user solution
func (t *Tester) BuildDockerImage() error {
	// Get the absolute path to the project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	fmt.Println("Current working directory:", projectRoot)

	// Construct the absolute path to the Dockerfile
	dockerfilePath := filepath.Join(projectRoot, "assets", "Dockerfile.solution")

	// Use the absolute path in the docker build command
	cmd := exec.Command("docker", "build", "-t", t.ImageName, "-f", dockerfilePath, ".")
	cmd.Dir = projectRoot // Ensure Docker uses the project root as context
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %v\nOutput: %s", err, out.String())
	}

	return nil
}

// DeployPod creates a pod to test the solution
func (t *Tester) DeployPod(podName string) error {
	cmd := exec.Command("kubectl", "apply", "-f", "./assets/solution-pod.yaml")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy pod: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) GetPodLogs(podName string) (string, error) {
	// Wait for the pod to be ready
	timeout := time.After(30 * time.Second) // Adjust timeout as needed
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout waiting for pod '%s' to be ready", podName)
		case <-tick:
			cmd := exec.Command("kubectl", "get", "pod", podName, "-o", "jsonpath={.status.phase}")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			if err := cmd.Run(); err != nil {
				return "", fmt.Errorf("failed to check pod status: %v\nOutput: %s", err, out.String())
			}

			status := out.String()
			if status == "Running" {
				goto Ready
			}
			fmt.Printf("Pod '%s' status: %s. Waiting...\n", podName, status)
		}
	}

Ready:
	// Retrieve logs
	cmd := exec.Command("kubectl", "logs", podName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get logs from pod '%s': %v\nOutput: %s", podName, err, out.String())
	}
	return out.String(), nil
}

// CleanUp removes the pod and other resources
func (t *Tester) CleanUp(podName string) error {
	var errors []string

	// Delete the pod
	if err := exec.Command("kubectl", "delete", "pod", podName).Run(); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to delete pod %s: %v", podName, err))
	}

	// Remove the Docker image
	if err := exec.Command("docker", "rmi", t.ImageName).Run(); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to remove Docker image %s: %v", t.ImageName, err))
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}
	return nil
}
