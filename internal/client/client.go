package client

import (
	"context"
	"fmt"

	"github.com/somaz94/kube-events/internal/event"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// EventLister defines the interface for listing Kubernetes events.
type EventLister interface {
	ListEvents(ctx context.Context, namespace string) ([]event.Event, error)
}

// Client wraps the Kubernetes clientset for event operations.
type Client struct {
	clientset kubernetes.Interface
}

// NewFromClientset creates a Client from an existing kubernetes.Interface (for testing).
func NewFromClientset(cs kubernetes.Interface) *Client {
	return &Client{clientset: cs}
}

// New creates a new Client using the given kubeconfig and context.
func New(kubeconfig, kubeContext string) (*Client, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}

	overrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{clientset: cs}, nil
}

// ListEvents fetches events from the specified namespace (empty string = all namespaces).
func (c *Client) ListEvents(ctx context.Context, namespace string) ([]event.Event, error) {
	list, err := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	events := make([]event.Event, 0, len(list.Items))
	for _, e := range list.Items {
		events = append(events, event.ConvertK8sEvent(e))
	}
	return events, nil
}
