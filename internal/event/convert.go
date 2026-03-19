package event

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// ConvertK8sEvent converts a Kubernetes corev1.Event to an internal Event.
func ConvertK8sEvent(e corev1.Event) Event {
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

	return Event{
		Type:      e.Type,
		Reason:    e.Reason,
		Message:   e.Message,
		Count:     e.Count,
		FirstSeen: firstSeen,
		LastSeen:  lastSeen,
		Age:       time.Since(lastSeen),
		InvolvedObject: InvolvedObject{
			Kind:      e.InvolvedObject.Kind,
			Name:      e.InvolvedObject.Name,
			Namespace: e.InvolvedObject.Namespace,
		},
		Source: Source{
			Component: e.Source.Component,
			Host:      e.Source.Host,
		},
	}
}
