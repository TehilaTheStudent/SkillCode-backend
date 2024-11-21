package tester

import (
	"fmt"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	language model.PredefinedSupportedLanguage // Language
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
func NewUniqueTester(sharedTester *SharedTester, jobName, imageName, runtimeCommand, fileExtension, requestID string,lanugage model.PredefinedSupportedLanguage ) *UniqueTester {
	return &UniqueTester{
		sharedTester:   sharedTester,
		jobName:        jobName,
		imageName:      imageName,
		runtimeCommand: runtimeCommand,
		fileExtension:  fileExtension,
		requestID:      requestID,
		configMapName:  fmt.Sprintf("user-script-%s", requestID),
		language: lanugage,
	}
}
