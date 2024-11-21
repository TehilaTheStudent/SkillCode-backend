package tester

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"time"

// 	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
// 	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
// 	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
// )

// func TestUserSolution(question *model.Question, userFunction string, language model.PredefinedSupportedLanguage, cfg config.ConfigSandbox) (*model.Feedback, error) {
// 	tester := NewTester()

// 	// Development mode: Simulate local execution
// 	var rawLogs string
// 	var err error
// 	if config.GlobalConfigAPI.ModeEnv == "development" {
// 		rawLogs, err = tester.executeLocally(language, cfg)
// 	} else {
// 		rawLogs, err = tester.executeWithDockerAndK8s(cfg)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parse JSON logs into Feedback struct
// 	var feedback model.Feedback
// 	if parseErr := json.Unmarshal([]byte(rawLogs), &feedback); parseErr != nil {
// 		return nil, model.NewCustomError(500, fmt.Sprintf("failed to parse feedback logs: %v", parseErr))
// 	}

// 	return &feedback, nil
// }
// func (t *Tester) executeLocally(language model.PredefinedSupportedLanguage, cfg config.ConfigSandbox) (string, error) {
// 	runTime := model.GetRuntime(language)
// 	var cmdArgs []string
// 	switch language {
// 	case model.Python:
// 		cmdArgs = []string{runTime, cfg.TestUserCodePath}
// 	case model.JavaScript:
// 		cmdArgs = []string{runTime, cfg.TestUserCodePath}
// 	default:
// 		return "", model.NewCustomError(400, "unsupported language")
// 	}

// 	// Capture the command's standard output and error
// 	stdout, err := utils.RunCommand(cmdArgs[0], cmdArgs[1:]...)
// 	if err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to execute local file: %v", err))
// 	}

// 	return stdout, nil
// }

// // Encapsulate the existing Docker and Kubernetes logic into a helper function
// func (t *Tester) executeWithDockerAndK8s(cfg config.LanguageConfig) (string, error) {
// 	defer func() {
// 		if err := t.CleanUp(cfg.PodName, cfg.ImageName); err != nil {
// 			fmt.Println("Failed to clean up resources:", err)
// 		}
// 	}()

// 	// 1. Build Docker image
// 	if err := t.BuildDockerImage(cfg.ImageName, cfg.DockerFilePath); err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to build Docker image: %v", err))
// 	}

// 	// 2. Load Docker image into Kind
// 	if err := t.LoadImageIntoKind(cfg.ImageName, cfg.ClusterName); err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to load image into Kind: %v", err))
// 	}

// 	// 3. Deploy the pod
// 	if err := t.DeployPod(cfg.PodName, cfg.PodFilePath); err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to deploy pod: %v", err))
// 	}

// 	// 4. Retrieve pod logs
// 	logs, err := t.GetPodLogs(cfg.PodName)
// 	if err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to get pod logs: %v", err))
// 	}

// 	return logs, nil
// }

// type Tester struct {
// 	// i have to add here fields, should it be for concurrent users later?
// }

// // DeployPod creates a pod to test the solution
// func (t *Tester) DeployPod(podName string, podPath string) error {
// 	_, err := utils.RunCommand("kubectl", "apply", "-f", podPath)
// 	if err != nil {
// 		return model.NewCustomError(500, fmt.Sprintf("failed to deploy pod: %v", err))
// 	}
// 	return nil
// }

// func (t *Tester) LoadImageIntoKind(imageName, clusterName string) error {
// 	_, err := utils.RunCommand("kind", "load", "docker-image", imageName, "--name", clusterName)
// 	if err != nil {
// 		return model.NewCustomError(500, fmt.Sprintf("failed to load image into Kind: %v", err))
// 	}
// 	return nil
// }

// func (t *Tester) SavePodLogs(podName, logFile string) error {
// 	logs, err := t.GetPodLogs(podName)
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(logFile, []byte(logs), 0644)
// }

// // BuildDockerImage builds a Docker image for the user solution
// func (t *Tester) BuildDockerImage(ImageName string, dockerPath string) error {
// 	// Get the absolute path to the project root
// 	projectRoot, err := os.Getwd()
// 	if err != nil {
// 		return model.NewCustomError(500, fmt.Sprintf("failed to get working directory: %v", err))
// 	}

// 	// Use the absolute path in the docker build command
// 	cmd := exec.Command("docker", "build", "-t", ImageName, "-f", dockerPath, ".")
// 	cmd.Dir = projectRoot // Ensure Docker uses the project root as context
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &out
// 	if err := cmd.Run(); err != nil {
// 		return model.NewCustomError(500, fmt.Sprintf("failed to build Docker image: %v\nOutput: %s", err, out.String()))
// 	}

// 	return nil
// }
// func NewTester() *Tester {
// 	return &Tester{}
// }

// func (t *Tester) GetPodLogs(podName string) (string, error) {
// 	// Wait for the pod to be ready
// 	timeout := time.After(30 * time.Second) // Adjust timeout as needed
// 	tick := time.Tick(2 * time.Second)

// 	for {
// 		select {
// 		case <-timeout:
// 			return "", model.NewCustomError(500, fmt.Sprintf("timeout waiting for pod '%s' to be ready", podName))
// 		case <-tick:
// 			out, err := utils.RunCommand("kubectl", "get", "pod", podName, "-o", "jsonpath={.status.phase}")
// 			if err != nil {
// 				return "", model.NewCustomError(500, fmt.Sprintf("failed to check pod status: %v", err))
// 			}

// 			status := out
// 			if status == "Running" {
// 				goto Ready
// 			}
// 			fmt.Printf("Pod '%s' status: %s. Waiting...\n", podName, status)
// 		}
// 	}

// Ready:
// 	// Retrieve logs
// 	out, err := utils.RunCommand("kubectl", "logs", podName)
// 	if err != nil {
// 		return "", model.NewCustomError(500, fmt.Sprintf("failed to get logs from pod '%s': %v", podName, err))
// 	}
// 	return out, nil
// }

// // CleanUp removes the pod and other resources
// func (t *Tester) CleanUp(podName string, imageName string) error {
// 	var errors []string

// 	// Step 1: Delete the pod
// 	cmd := exec.Command("kubectl", "delete", "pod", podName, "--ignore-not-found")
// 	var podOut bytes.Buffer
// 	cmd.Stdout = &podOut
// 	cmd.Stderr = &podOut
// 	if err := cmd.Run(); err != nil {
// 		errors = append(errors, fmt.Sprintf("Failed to delete pod %s: %v\nOutput: %s", podName, err, podOut.String()))
// 	}

// 	// Step 2: Remove the Docker image
// 	cmd = exec.Command("docker", "rmi", imageName)
// 	var imageOut bytes.Buffer
// 	cmd.Stdout = &imageOut
// 	cmd.Stderr = &imageOut
// 	if err := cmd.Run(); err != nil && !strings.Contains(imageOut.String(), "No such image") {
// 		errors = append(errors, fmt.Sprintf("Failed to remove Docker image %s: %v\nOutput: %s", imageName, err, imageOut.String()))
// 	}

// 	// Step 4: Check if there are any accumulated errors
// 	if len(errors) > 0 {
// 		return model.NewCustomError(500, strings.Join(errors, "; "))
// 	}
// 	return nil
// }
