package client

import (
	"context"
	"fmt"
	"time"

	"github.com/somaz94/kube-events/internal/event"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes clientset for event operations.
type Client struct {
	clientset kubernetes.Interface
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
		events = append(events, convertEvent(e))
	}
	return events, nil
}

func convertEvent(e corev1.Event) event.Event {
	lastSeen := e.LastTimestamp.Time
	if lastSeen.IsZero() {
		lastSeen = e.EventTime.Time
	}
	if lastSeen.IsZero() {
		lastSeen = e.CreationTimestamp.Time
	}

	firstSeen := e.FirstTimestamp.Time
	if firstSeen.IsZero() {
		firstSeen = lastSeen
	}

	return event.Event{
		Type:      e.Type,
		Reason:    e.Reason,
		Message:   e.Message,
		Count:     e.Count,
		FirstSeen: firstSeen,
		LastSeen:  lastSeen,
		Age:       time.Since(lastSeen),
		InvolvedObject: event.InvolvedObject{
			Kind:      e.InvolvedObject.Kind,
			Name:      e.InvolvedObject.Name,
			Namespace: e.InvolvedObject.Namespace,
		},
		Source: event.Source{
			Component: e.Source.Component,
			Host:      e.Source.Host,
		},
	}
}
