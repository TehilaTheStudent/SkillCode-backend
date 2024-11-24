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

	"sigs.k8s.io/yaml"
)





// ExecuteWithJobTemplate runs a user script via a Kubernetes Job
func (t *UniqueTester) ExecuteWithJobTemplate(params map[string]string, jobTemplatePath, scriptContent string) (string, error) {
	// Step 1: Create the ConfigMap
	if err := t.CreateConfigMap(scriptContent); err != nil {
		return "", fmt.Errorf("failed to create ConfigMap: %v", err)
	}

	// Print the ConfigMap content for debugging
	// fmt.Printf("ConfigMap '%s' content:\n%s\n", t.configMapName, scriptContent)

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

	// Print the processed job for debugging
	// fmt.Printf("Processed Job YAML:\n%s\n", processedJob.String())

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
    timeout := time.After(30 * time.Second) // Adjusted timeout to 30 seconds
    tick := time.Tick(1 * time.Second)     // Reduced tick interval for quicker checks

    for {
        select {
        case <-timeout:
            return "", fmt.Errorf("timeout waiting for Job '%s' to complete after 30 seconds", jobName)
        case <-tick:
            // Fetch the Job object
            job, err := t.sharedTester.ClientSet.BatchV1().Jobs(t.sharedTester.Namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
            if err != nil {
                return "", fmt.Errorf("failed to get Job status: %v", err)
            }

            // Check if the Job succeeded
            if job.Status.Succeeded > 0 {
                return t.getJobLogs(jobName) // Fetch logs immediately if the job succeeded
            }

            // Check if the Job failed
            if job.Status.Failed > 0 {
                logs, logErr := t.getJobLogs(jobName)
                if logErr != nil {
                    return "", fmt.Errorf("Job '%s' failed and failed to get logs: %v", jobName, logErr)
                }
                return "", fmt.Errorf("Job '%s' failed. Logs:\n%s", jobName, logs)
            }

            // Optional: Log progress for debugging
            fmt.Printf("Job '%s' is still running...\n", jobName)
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
