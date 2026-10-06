package client

import (
	"context"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListEvents(t *testing.T) {
	now := time.Now()
	fakeCS := fake.NewSimpleClientset(
		&corev1.Event{
			ObjectMeta:     metav1.ObjectMeta{Name: "evt-1", Namespace: "default"},
			Type:           "Warning",
			Reason:         "BackOff",
			Message:        "Back-off restarting failed container",
			Count:          3,
			LastTimestamp:  metav1.Time{Time: now.Add(-5 * time.Minute)},
			FirstTimestamp: metav1.Time{Time: now.Add(-10 * time.Minute)},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "app-1",
				Namespace: "default",
			},
			Source: corev1.EventSource{Component: "kubelet", Host: "node-1"},
		},
		&corev1.Event{
			ObjectMeta:     metav1.ObjectMeta{Name: "evt-2", Namespace: "default"},
			Type:           "Normal",
			Reason:         "Scheduled",
			Message:        "Successfully assigned default/app-1 to node-1",
			Count:          1,
			LastTimestamp:  metav1.Time{Time: now.Add(-8 * time.Minute)},
			FirstTimestamp: metav1.Time{Time: now.Add(-8 * time.Minute)},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "app-1",
				Namespace: "default",
			},
		},
		&corev1.Event{
			ObjectMeta:    metav1.ObjectMeta{Name: "evt-3", Namespace: "production"},
			Type:          "Normal",
			Reason:        "Pulling",
			Message:       "Pulling image nginx:1.27",
			Count:         1,
			LastTimestamp: metav1.Time{Time: now.Add(-2 * time.Minute)},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "web-1",
				Namespace: "production",
			},
		},
	)

	c := NewFromClientset(fakeCS)
	ctx := context.Background()

	// List all events
	events, err := c.ListEvents(ctx, "")
	if err != nil {
		t.Fatalf("ListEvents error: %v", err)
	}
	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}

	events, err = c.ListEvents(ctx, "default")
	if err != nil {
		t.Fatalf("ListEvents(default) error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events in default, got %d", len(events))
	}

	events, err = c.ListEvents(ctx, "production")
	if err != nil {
		t.Fatalf("ListEvents(production) error: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event in production, got %d", len(events))
	}
}

func TestListEvents_Empty(t *testing.T) {
	fakeCS := fake.NewSimpleClientset()
	c := NewFromClientset(fakeCS)

	events, err := c.ListEvents(context.Background(), "")
	if err != nil {
		t.Fatalf("ListEvents error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestNewFromClientset(t *testing.T) {
	fakeCS := fake.NewSimpleClientset()
	c := NewFromClientset(fakeCS)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.clientset != fakeCS {
		t.Error("expected clientset to match")
	}
}

func TestNew_InvalidKubeconfig(t *testing.T) {
	_, err := New("/nonexistent/kubeconfig", "")
	if err == nil {
		t.Error("expected error for invalid kubeconfig path")
	}
}

func TestNew_InvalidContext(t *testing.T) {
	tmpFile := t.TempDir() + "/kubeconfig"
	content := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://localhost:6443
  name: test
contexts:
- context:
    cluster: test
    user: test
  name: test
current-context: test
users:
- name: test
  user:
    token: fake-token
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := New(tmpFile, "nonexistent-context")
	if err == nil {
		t.Error("expected error for nonexistent context")
	}
}

func TestNew_ValidKubeconfig(t *testing.T) {
	tmpFile := t.TempDir() + "/kubeconfig"
	content := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://localhost:6443
  name: test
contexts:
- context:
    cluster: test
    user: test
  name: test
current-context: test
users:
- name: test
  user:
    token: fake-token
`
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := New(tmpFile, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Error("expected non-nil client")
	}
}
