package k8s

import (
	"context"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
	namespace string
	timeout   time.Duration
}

func NewK8sClient(namespace string) (*K8sClient, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		clientset: cs,
		namespace: namespace,
		timeout:   30 * time.Second,
	}, nil
}

func (c *K8sClient) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.timeout)
}

// Helpers to access typed clients
func (c *K8sClient) JobsClient() kubernetes.Interface {
	return c.clientset
}

func (c *K8sClient) Namespace() string {
	return c.namespace
}

// Convenience wrappers
func (c *K8sClient) CreateJob(ctx context.Context, job *batchv1.Job) (*batchv1.Job, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.clientset.BatchV1().Jobs(c.namespace).Create(ctx, job, metav1.CreateOptions{})
}

func (c *K8sClient) GetJob(ctx context.Context, name string) (*batchv1.Job, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.clientset.BatchV1().Jobs(c.namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *K8sClient) DeleteJob(ctx context.Context, name string) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	propagation := metav1.DeletePropagationForeground
	return c.clientset.BatchV1().Jobs(c.namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
}

func (c *K8sClient) ListPodsForJob(ctx context.Context, jobName string) (*corev1.PodList, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.clientset.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	})
}
