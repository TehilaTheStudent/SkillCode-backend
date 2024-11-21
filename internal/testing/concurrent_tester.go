package tester

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	v1 "k8s.io/api/batch/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// SharedTester handles shared Kubernetes resources
type SharedTester struct {
	ClientSet *kubernetes.Clientset // Kubernetes client, shared among all testers.
	Namespace string                // Namespace where the Job will be created.
}

// UniqueTester holds request-specific data and uses SharedTester for shared resources
type UniqueTester struct {
	sharedTester   *SharedTester // Reference to SharedTester
	jobName        string        // Unique Job name for the request.
	imageName      string        // Image name used for the Job.
	runtimeCommand string        // Runtime command for the user's script.
	fileExtension  string        // File extension for the script (e.g., .py, .js).
	requestID      string        // Unique request identifier.
	configMapName  string        // ConfigMap name for the user's script.
}

// NewSharedTester initializes SharedTester with Kubernetes client
func NewSharedTester(kubeconfigPath, namespace string) (*SharedTester, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	ClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return &SharedTester{
		ClientSet: ClientSet,
		Namespace: namespace,
	}, nil
}

// NewUniqueTester initializes UniqueTester with request-specific data
func NewUniqueTester(sharedTester *SharedTester, jobName, imageName, runtimeCommand, fileExtension, requestID string) *UniqueTester {
	return &UniqueTester{
		sharedTester:   sharedTester,
		jobName:        jobName,
		imageName:      imageName,
		runtimeCommand: runtimeCommand,
		fileExtension:  fileExtension,
		requestID:      requestID,
		configMapName:  fmt.Sprintf("user-script-%s", requestID),
	}
}

// ExecuteWithJobTemplate runs a user script via a Kubernetes Job
func (t *UniqueTester) ExecuteWithJobTemplate(params map[string]string, jobTemplatePath, scriptContent string) (string, error) {
	// Step 1: Create the ConfigMap
	if err := t.CreateConfigMap(scriptContent); err != nil {
		return "", fmt.Errorf("failed to create ConfigMap: %v", err)
	}

	// Ensure ConfigMap cleanup after execution
	defer func() {
		if err := t.DeleteConfigMap(); err != nil {
			fmt.Printf("Warning: failed to clean up ConfigMap: %v\n", err)
		}
	}()

	// Step 2: Load the Job template
	templateBytes, err := os.ReadFile(jobTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read Job template: %v", err)
	}

	// Step 3: Process the template
	tmpl, err := template.New("job").Option("missingkey=error").Parse(string(templateBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse Job template: %v", err)
	}

	var processedJob bytes.Buffer
	if err := tmpl.Execute(&processedJob, params); err != nil {
		return "", fmt.Errorf("failed to execute Job template: %v", err)
	}

	// Step 4: Convert the processed YAML into a Job object
	job := &v1.Job{}
	if err := yaml.Unmarshal(processedJob.Bytes(), job); err != nil {
		return "", fmt.Errorf("failed to unmarshal Job YAML: %v", err)
	}

	// Step 5: Submit the Job to Kubernetes
	job, err = t.sharedTester.ClientSet.BatchV1().Jobs(t.sharedTester.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create Job: %v", err)
	}

	// Step 6: Wait for Job completion and retrieve logs
	return t.waitForJobAndFetchLogs(job.Name)
}

// waitForJobAndFetchLogs waits for the Job to complete (success or failure) and fetches logs if successful.
func (t *UniqueTester) waitForJobAndFetchLogs(jobName string) (string, error) {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout waiting for Job '%s' to complete", jobName)
		case <-tick:
			// Fetch the Job object
			job, err := t.sharedTester.ClientSet.BatchV1().Jobs(t.sharedTester.Namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
			if err != nil {
				return "", fmt.Errorf("failed to get Job status: %v", err)
			}

			// Check if the Job succeeded
			if job.Status.Succeeded > 0 {
				// Retrieve logs from the Job's pods
				return t.getJobLogs(jobName)
			}

			// Check if the Job failed
			if job.Status.Failed > 0 {
				return "", fmt.Errorf("Job '%s' failed. Check logs or events for more details.", jobName)
			}

			// Check if Job conditions indicate failure
			for _, condition := range job.Status.Conditions {
				if condition.Type == "Failed" && condition.Status == "True" {
					return "", fmt.Errorf("Job '%s' failed with reason: %s, message: %s", jobName, condition.Reason, condition.Message)
				}
			}

			// Log Job progress
			fmt.Printf("Job '%s' is still running. Waiting...\n", jobName)
		}
	}
}


// getJobLogs retrieves logs from the Pod associated with the Job
func (t *UniqueTester) getJobLogs(jobName string) (string, error) {
	pods, err := t.sharedTester.ClientSet.CoreV1().Pods(t.sharedTester.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list Pods for Job '%s': %v", jobName, err)
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no Pods found for Job '%s'", jobName)
	}

	podName := pods.Items[0].Name
	logs, err := t.sharedTester.ClientSet.CoreV1().
		Pods(t.sharedTester.Namespace).
		GetLogs(podName, &corev1.PodLogOptions{}).
		DoRaw(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to get logs for Pod '%s': %v", podName, err)
	}

	return string(logs), nil
}

// CreateConfigMap creates a ConfigMap for the user script
func (t *UniqueTester) CreateConfigMap(scriptContent string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.configMapName,
			Namespace: t.sharedTester.Namespace,
		},
		Data: map[string]string{
			"run_tests": scriptContent, // Key in the ConfigMap
		},
	}

	_, err := t.sharedTester.ClientSet.CoreV1().ConfigMaps(t.sharedTester.Namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ConfigMap: %v", err)
	}

	fmt.Printf("ConfigMap '%s' created successfully\n", t.configMapName)
	return nil
}

// DeleteConfigMap deletes the ConfigMap after the Job is finished
func (t *UniqueTester) DeleteConfigMap() error {
	err := t.sharedTester.ClientSet.CoreV1().ConfigMaps(t.sharedTester.Namespace).Delete(context.TODO(), t.configMapName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete ConfigMap '%s': %v", t.configMapName, err)
	}

	fmt.Printf("ConfigMap '%s' deleted successfully\n", t.configMapName)
	return nil
}
