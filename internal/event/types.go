package event

import "time"

// Event represents a Kubernetes event with relevant fields.
type Event struct {
	Type           string
	Reason         string
	Message        string
	Count          int32
	FirstSeen      time.Time
	LastSeen       time.Time
	Age            time.Duration
	InvolvedObject InvolvedObject
	Source         Source
}

// InvolvedObject identifies the resource that the event is about.
type InvolvedObject struct {
	Kind      string
	Name      string
	Namespace string
}

// Source identifies the component that reported the event.
type Source struct {
	Component string
	Host      string
}

// ResourceKey uniquely identifies a resource involved in events.
type ResourceKey struct {
	Kind      string
	Name      string
	Namespace string
}

// ResourceGroup holds events grouped by their involved resource.
type ResourceGroup struct {
	Key    ResourceKey
	Events []Event
}
