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

// Config holds the dynamic configuration values
type Config struct {
	ImageName     string
	DockerFilePath string
	PodFilePath   string
	PodName       string
	ClusterName   string
}

func NewConfig() *Config {
	return &Config{
		ImageName:     "user-solution:latest",
		DockerFilePath: "./assets/Dockerfile.python",
		PodFilePath:   "./assets/solution_pod.yaml",
		PodName:       "solution-pod",
		ClusterName:   "my-cluster",
	}
}

// executeFunction is a placeholder for the logic to execute a user's function with test cases.
// This function should be implemented in a real environment with sandboxing.
func TestUserSolution(question *model.Question, userFunction string, language string) (string, error) {
	tester := NewTester()
	config := NewConfig()
	// Ensure cleanup always happens, no matter where the function exits
	defer func() {
		if err := tester.CleanUp(config.PodName, config.ImageName); err != nil {
			fmt.Printf("Cleanup failed: %v\n", err)
		}
	}()

	if err := tester.EnsureKubectlInstalled(); err != nil {
		return "", utils.New(500, fmt.Sprintf("kubectl check failed: %v", err))
	}

	if err := tester.EnsureClusterExists(config.ClusterName); err != nil {
		return "", utils.New(500, fmt.Sprintf("failed to create Kind cluster: %v", err))
	}

	//  Build Docker image
	if err := tester.BuildDockerImage(config.ImageName, config.DockerFilePath); err != nil {
		return "", utils.New(500, fmt.Sprintf("failed to build Docker image: %v", err))
	}
	//  Load Docker image into Kind
	if err := tester.LoadImageIntoKind(config.ImageName, config.ClusterName); err != nil {
		return "", utils.New(500, fmt.Sprintf("failed to load image: %v", err))
	}
	//  Set kubectl context
	if err := tester.EnsureKindContext(config.ClusterName); err != nil {
		return "", utils.New(500, fmt.Sprintf("failed to set kubectl context: %v", err))
	}

	//  Deploy the pod
	if err := tester.DeployPod(config.PodName, config.PodFilePath); err != nil {
		return "", utils.New(500, "failed to deploy pod: "+err.Error())
	}

	//  Retrieve pod logs
	logs, err := tester.GetPodLogs(config.PodName)
	if err != nil {
		return "", utils.New(500, "failed to get pod logs: "+err.Error())
	}
	// Save logs to a file
	// if err := tester.SavePodLogs(config.PodName, "./temp/python/logs.txt"); err != nil {
	// 	return "", utils.New(500, "failed to save pod logs: "+err.Error())
	// }

	return logs, nil
}

type Tester struct {
	// add here
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", utils.New(500, fmt.Sprintf("command failed: %v\nOutput: %s", err, out.String()))
	}
	return out.String(), nil
}

func (t *Tester) EnsureKubectlInstalled() error {
	_, err := runCommand("kubectl", "version", "--client")
	if err != nil {
		return utils.New(500, fmt.Sprintf("kubectl not installed or misconfigured: %v", err))
	}
	return nil
}

// DeployPod creates a pod to test the solution
func (t *Tester) DeployPod(podName string, podPath string) error {
	_, err := runCommand("kubectl", "apply", "-f", podPath)
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to deploy pod: %v", err))
	}
	return nil
}

func (t *Tester) EnsureClusterExists(clusterName string) error {
	out, err := runCommand("kind", "get", "clusters")
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to get Kind clusters: %v", err))
	}
	clusters := strings.Split(out, "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return nil // Cluster exists
		}
	}
	_, err = runCommand("kind", "create", "cluster", "--name", clusterName)
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to create Kind cluster: %v", err))
	}
	return nil
}

func (t *Tester) LoadImageIntoKind(imageName, clusterName string) error {
	_, err := runCommand("kind", "load", "docker-image", imageName, "--name", clusterName)
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to load image into Kind: %v", err))
	}
	return nil
}

func (t *Tester) EnsureKindContext(clusterName string) error {
	_, err := runCommand("kubectl", "config", "use-context", fmt.Sprintf("kind-%s", clusterName))
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to switch kubectl context: %v", err))
	}
	return nil
}

func (t *Tester) DeleteKindCluster(clusterName string) error {
	_, err := runCommand("kind", "delete", "cluster", "--name", clusterName)
	if err != nil {
		return utils.New(500, fmt.Sprintf("failed to delete Kind cluster: %v", err))
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
		return utils.New(500, fmt.Sprintf("failed to get working directory: %v", err))
	}

	// Use the absolute path in the docker build command
	cmd := exec.Command("docker", "build", "-t", ImageName, "-f", dockerPath, ".")
	cmd.Dir = projectRoot // Ensure Docker uses the project root as context
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return utils.New(500, fmt.Sprintf("failed to build Docker image: %v\nOutput: %s", err, out.String()))
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
			return "", utils.New(500, fmt.Sprintf("timeout waiting for pod '%s' to be ready", podName))
		case <-tick:
			out, err := runCommand("kubectl", "get", "pod", podName, "-o", "jsonpath={.status.phase}")
			if err != nil {
				return "", utils.New(500, fmt.Sprintf("failed to check pod status: %v", err))
			}

			status := out
			if status == "Running" {
				goto Ready
			}
			fmt.Printf("Pod '%s' status: %s. Waiting...\n", podName, status)
		}
	}

Ready:
	// Retrieve logs
	out, err := runCommand("kubectl", "logs", podName)
	if err != nil {
		return "", utils.New(500, fmt.Sprintf("failed to get logs from pod '%s': %v", podName, err))
	}
	return out, nil
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

	// Step 4: Check if there are any accumulated errors
	if len(errors) > 0 {
		return utils.New(500, strings.Join(errors, "; "))
	}
	return nil
}
