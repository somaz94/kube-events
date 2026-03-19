package event

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertK8sEvent_AllFields(t *testing.T) {
	now := time.Now()
	k8sEvent := corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "evt-1", Namespace: "default"},
		Type:       "Warning",
		Reason:     "BackOff",
		Message:    "Back-off restarting",
		Count:      3,
		LastTimestamp:  metav1.Time{Time: now.Add(-2 * time.Minute)},
		FirstTimestamp: metav1.Time{Time: now.Add(-5 * time.Minute)},
		InvolvedObject: corev1.ObjectReference{
			Kind: "Pod", Name: "app-1", Namespace: "default",
		},
		Source: corev1.EventSource{Component: "kubelet", Host: "node-1"},
	}

	e := ConvertK8sEvent(k8sEvent)

	if e.Type != "Warning" {
		t.Errorf("Type = %q, want Warning", e.Type)
	}
	if e.Reason != "BackOff" {
		t.Errorf("Reason = %q, want BackOff", e.Reason)
	}
	if e.Message != "Back-off restarting" {
		t.Errorf("Message = %q, want 'Back-off restarting'", e.Message)
	}
	if e.Count != 3 {
		t.Errorf("Count = %d, want 3", e.Count)
	}
	if e.InvolvedObject.Kind != "Pod" {
		t.Errorf("Kind = %q, want Pod", e.InvolvedObject.Kind)
	}
	if e.InvolvedObject.Name != "app-1" {
		t.Errorf("Name = %q, want app-1", e.InvolvedObject.Name)
	}
	if e.InvolvedObject.Namespace != "default" {
		t.Errorf("Namespace = %q, want default", e.InvolvedObject.Namespace)
	}
	if e.Source.Component != "kubelet" {
		t.Errorf("Component = %q, want kubelet", e.Source.Component)
	}
	if e.Source.Host != "node-1" {
		t.Errorf("Host = %q, want node-1", e.Source.Host)
	}
	if e.LastSeen.IsZero() {
		t.Error("LastSeen should not be zero")
	}
	if e.FirstSeen.IsZero() {
		t.Error("FirstSeen should not be zero")
	}
	if e.Age <= 0 {
		t.Error("Age should be positive")
	}
}

func TestConvertK8sEvent_EventTimeFallback(t *testing.T) {
	now := time.Now()
	e := ConvertK8sEvent(corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "e1"},
		EventTime:  metav1.MicroTime{Time: now.Add(-1 * time.Minute)},
	})
	if e.LastSeen.IsZero() {
		t.Error("expected LastSeen from EventTime")
	}
}

func TestConvertK8sEvent_CreationTimestampFallback(t *testing.T) {
	now := time.Now()
	e := ConvertK8sEvent(corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "e2",
			CreationTimestamp: metav1.Time{Time: now},
		},
	})
	if e.LastSeen.IsZero() {
		t.Error("expected LastSeen from CreationTimestamp")
	}
}

func TestConvertK8sEvent_FirstSeenFallback(t *testing.T) {
	now := time.Now()
	e := ConvertK8sEvent(corev1.Event{
		ObjectMeta:    metav1.ObjectMeta{Name: "e3"},
		LastTimestamp:  metav1.Time{Time: now.Add(-3 * time.Minute)},
	})
	if e.FirstSeen.IsZero() {
		t.Error("expected FirstSeen fallback to LastSeen")
	}
	if !e.FirstSeen.Equal(e.LastSeen) {
		t.Error("FirstSeen should equal LastSeen when FirstTimestamp is zero")
	}
}
