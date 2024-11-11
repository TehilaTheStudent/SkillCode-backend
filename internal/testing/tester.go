package tester

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
)

// executeFunction is a placeholder for the logic to execute a user's function with test cases.
// This function should be implemented in a real environment with sandboxing.
func TestUserSolution(question *model.Question, userFunction string, language string) (string, error) {
	tester := NewTester()
	// Ensure cleanup always happens, no matter where the function exits
	imageName := "user-solution:latest"
	dockerFilePath := "./assets/Dockerfile.python"
	podFilePath := "./assets/solution_pod.yaml"
	podName := "solution-pod"
	defer func() {
		if err := tester.CleanUp(podName, imageName); err != nil {
			fmt.Printf("Cleanup failed: %v\n", err)
		}
	}()

	if err := tester.EnsureKubectlInstalled(); err != nil {
		return "", fmt.Errorf("kubectl check failed: %w", err)
	}

	// Step 1: Check if the language is supported
	// var functionSignature string
	// found := false
	// for _, langConfig := range question.Languages {
	// 	if langConfig.Language == language {
	// 		functionSignature = langConfig.FunctionSignature
	// 		found = true
	// 		break
	// 	}
	// }
	// if !found {
	// 	return "", utils.New(404, "language "+language+" not supported")
	// }

	// if err := generatePythonScript("./temp/python/main.py", userFunction, functionSignature, question.TestCases); err != nil {
	// 	return "", utils.New(500, "failed to generate main file: "+err.Error())
	// }

	if err := tester.EnsureClusterExists("my-cluster"); err != nil {
		return "", fmt.Errorf("failed to create Kind cluster: %w", err)
	}

	//  Build Docker image
	if err := tester.BuildDockerImage(imageName, dockerFilePath); err != nil {
		return "", fmt.Errorf("failed to build Docker image: %w", err)
	}
	//  Load Docker image into Kind
	if err := tester.LoadImageIntoKind(imageName, "my-cluster"); err != nil {
		return "", fmt.Errorf("failed to load image: %w", err)
	}
	//  Set kubectl context
	if err := tester.EnsureKindContext("my-cluster"); err != nil {
		return "", fmt.Errorf("failed to set kubectl context: %w", err)
	}

	//  Deploy the pod
	if err := tester.DeployPod(podName, podFilePath); err != nil {
		return "", utils.New(500, "failed to deploy pod: "+err.Error())
	}

	//  Retrieve pod logs
	logs, err := tester.GetPodLogs(podName)
	if err != nil {
		return "", utils.New(500, "failed to get pod logs: "+err.Error())
	}
	// Save logs to a file
	// if err := tester.SavePodLogs(podName, "./temp/python/logs.txt"); err != nil {
	// 	return "", utils.New(500, "failed to save pod logs: "+err.Error())
	// }

	//  Clean up resources
	if err = tester.CleanUp(podName, imageName); err != nil {
		return "", utils.New(500, "failed to clean up resources: "+err.Error())
	}
	return logs, nil
}

type Tester struct {
	// add here
}

func (t *Tester) EnsureKubectlInstalled() error {
	cmd := exec.Command("kubectl", "version", "--client")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl not installed or misconfigured: %v\nOutput: %s", err, out.String())
	}
	return nil
}

// DeployPod creates a pod to test the solution
func (t *Tester) DeployPod(podName string, podPath string) error {
	cmd := exec.Command("kubectl", "apply", "-f", podPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy pod: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) EnsureClusterExists(clusterName string) error {
	cmd := exec.Command("kind", "get", "clusters")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get Kind clusters: %v\nOutput: %s", err, out.String())
	}
	clusters := strings.Split(out.String(), "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return nil // Cluster exists
		}
	}
	// Create the cluster if it doesn't exist
	cmd = exec.Command("kind", "create", "cluster", "--name", clusterName)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Kind cluster: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) LoadImageIntoKind(imageName, clusterName string) error {
	cmd := exec.Command("kind", "load", "docker-image", imageName, "--name", clusterName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load image into Kind: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) EnsureKindContext(clusterName string) error {
	cmd := exec.Command("kubectl", "config", "use-context", fmt.Sprintf("kind-%s", clusterName))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to switch kubectl context: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) DeleteKindCluster(clusterName string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete Kind cluster: %v\nOutput: %s", err, out.String())
	}
	return nil
}

func (t *Tester) SavePodLogs(podName, logFile string) error {
	logs, err := t.GetPodLogs(podName)
	if err != nil {
		return err
	}
	return os.WriteFile(logFile, []byte(logs), 0644)
}

// BuildDockerImage builds a Docker image for the user solution
func (t *Tester) BuildDockerImage(ImageName string, dockerPath string) error {
	// Get the absolute path to the project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Use the absolute path in the docker build command
	cmd := exec.Command("docker", "build", "-t", ImageName, "-f", dockerPath, ".")
	cmd.Dir = projectRoot // Ensure Docker uses the project root as context
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %v\nOutput: %s", err, out.String())
	}

	return nil
}
func NewTester() *Tester {
	return &Tester{}
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
func (t *Tester) CleanUp(podName string, imageName string) error {
	var errors []string

	// Step 1: Delete the pod
	cmd := exec.Command("kubectl", "delete", "pod", podName, "--ignore-not-found")
	var podOut bytes.Buffer
	cmd.Stdout = &podOut
	cmd.Stderr = &podOut
	if err := cmd.Run(); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to delete pod %s: %v\nOutput: %s", podName, err, podOut.String()))
	}

	// Step 2: Remove the Docker image
	cmd = exec.Command("docker", "rmi", imageName)
	var imageOut bytes.Buffer
	cmd.Stdout = &imageOut
	cmd.Stderr = &imageOut
	if err := cmd.Run(); err != nil && !strings.Contains(imageOut.String(), "No such image") {
		errors = append(errors, fmt.Sprintf("Failed to remove Docker image %s: %v\nOutput: %s", imageName, err, imageOut.String()))
	}

	// Step 3: Delete temporary log files
	// logFile := "./temp/python/logs.txt"
	// if _, err := os.Stat(logFile); err == nil {
	// 	if err := os.Remove(logFile); err != nil {
	// 		errors = append(errors, fmt.Sprintf("Failed to remove log file %s: %v", logFile, err))
	// 	}
	// }

	// Step 4: Check if there are any accumulated errors
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}
	return nil
}
