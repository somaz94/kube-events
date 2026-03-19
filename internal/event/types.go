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

// GroupBy defines how events are grouped in output.
type GroupBy string

const (
	GroupResource  GroupBy = "resource"
	GroupNamespace GroupBy = "namespace"
	GroupKind      GroupBy = "kind"
	GroupReason    GroupBy = "reason"
)

// ValidGroupBy returns true if the value is a supported group-by mode.
func ValidGroupBy(s string) bool {
	switch GroupBy(s) {
	case "", GroupResource, GroupNamespace, GroupKind, GroupReason:
		return true
	}
	return false
}

// ResourceKey uniquely identifies a resource involved in events.
type ResourceKey struct {
	Kind      string
	Name      string
	Namespace string
	Label     string // display label for non-resource grouping modes
}

// ResourceGroup holds events grouped by their involved resource.
type ResourceGroup struct {
	Key    ResourceKey
	Events []Event
}
