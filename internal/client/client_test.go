package client

import (
	"context"
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
			ObjectMeta: metav1.ObjectMeta{Name: "evt-1", Namespace: "default"},
			Type:       "Warning",
			Reason:     "BackOff",
			Message:    "Back-off restarting failed container",
			Count:      3,
			LastTimestamp: metav1.Time{Time: now.Add(-5 * time.Minute)},
			FirstTimestamp: metav1.Time{Time: now.Add(-10 * time.Minute)},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "app-1",
				Namespace: "default",
			},
			Source: corev1.EventSource{Component: "kubelet", Host: "node-1"},
		},
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{Name: "evt-2", Namespace: "default"},
			Type:       "Normal",
			Reason:     "Scheduled",
			Message:    "Successfully assigned default/app-1 to node-1",
			Count:      1,
			LastTimestamp: metav1.Time{Time: now.Add(-8 * time.Minute)},
			FirstTimestamp: metav1.Time{Time: now.Add(-8 * time.Minute)},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "app-1",
				Namespace: "default",
			},
		},
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{Name: "evt-3", Namespace: "production"},
			Type:       "Normal",
			Reason:     "Pulling",
			Message:    "Pulling image nginx:1.27",
			Count:      1,
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

	// List namespace-scoped events
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

func TestConvertEvent_Fields(t *testing.T) {
	now := time.Now()
	k8sEvent := corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "evt-1", Namespace: "default"},
		Type:       "Warning",
		Reason:     "Unhealthy",
		Message:    "Readiness probe failed",
		Count:      5,
		LastTimestamp:  metav1.Time{Time: now.Add(-3 * time.Minute)},
		FirstTimestamp: metav1.Time{Time: now.Add(-10 * time.Minute)},
		InvolvedObject: corev1.ObjectReference{
			Kind:      "Pod",
			Name:      "api-server",
			Namespace: "production",
		},
		Source: corev1.EventSource{Component: "kubelet", Host: "node-2"},
	}

	e := convertEvent(k8sEvent)

	if e.Type != "Warning" {
		t.Errorf("expected Type=Warning, got %s", e.Type)
	}
	if e.Reason != "Unhealthy" {
		t.Errorf("expected Reason=Unhealthy, got %s", e.Reason)
	}
	if e.Message != "Readiness probe failed" {
		t.Errorf("expected Message='Readiness probe failed', got %s", e.Message)
	}
	if e.Count != 5 {
		t.Errorf("expected Count=5, got %d", e.Count)
	}
	if e.InvolvedObject.Kind != "Pod" {
		t.Errorf("expected Kind=Pod, got %s", e.InvolvedObject.Kind)
	}
	if e.InvolvedObject.Name != "api-server" {
		t.Errorf("expected Name=api-server, got %s", e.InvolvedObject.Name)
	}
	if e.InvolvedObject.Namespace != "production" {
		t.Errorf("expected Namespace=production, got %s", e.InvolvedObject.Namespace)
	}
	if e.Source.Component != "kubelet" {
		t.Errorf("expected Source.Component=kubelet, got %s", e.Source.Component)
	}
	if e.Source.Host != "node-2" {
		t.Errorf("expected Source.Host=node-2, got %s", e.Source.Host)
	}
}

func TestConvertEvent_FallbackTimestamps(t *testing.T) {
	now := time.Now()

	// EventTime fallback (no LastTimestamp)
	e1 := convertEvent(corev1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: "e1"},
		EventTime:      metav1.MicroTime{Time: now.Add(-1 * time.Minute)},
		InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "p1"},
	})
	if e1.LastSeen.IsZero() {
		t.Error("expected LastSeen from EventTime, got zero")
	}

	// CreationTimestamp fallback (no LastTimestamp, no EventTime)
	e2 := convertEvent(corev1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: "e2", CreationTimestamp: metav1.Time{Time: now.Add(-2 * time.Minute)}},
		InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "p2"},
	})
	if e2.LastSeen.IsZero() {
		t.Error("expected LastSeen from CreationTimestamp, got zero")
	}

	// FirstSeen fallback to LastSeen
	e3 := convertEvent(corev1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: "e3"},
		LastTimestamp:   metav1.Time{Time: now.Add(-3 * time.Minute)},
		InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "p3"},
	})
	if e3.FirstSeen.IsZero() {
		t.Error("expected FirstSeen fallback to LastSeen, got zero")
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
