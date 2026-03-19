package event

import (
	"testing"
	"time"
)

func newEvent(typ, kind, name, ns, reason, msg string, age time.Duration) Event {
	return Event{
		Type:    typ,
		Reason:  reason,
		Message: msg,
		Count:   1,
		LastSeen:  time.Now().Add(-age),
		FirstSeen: time.Now().Add(-age),
		Age:       age,
		InvolvedObject: InvolvedObject{
			Kind:      kind,
			Name:      name,
			Namespace: ns,
		},
	}
}

func TestFilter_Since(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 30*time.Minute),
		newEvent("Normal", "Pod", "app-2", "default", "Scheduled", "scheduled", 2*time.Hour),
		newEvent("Warning", "Pod", "app-3", "default", "Unhealthy", "probe failed", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour})
	if len(result) != 2 {
		t.Errorf("expected 2 events within 1h, got %d", len(result))
	}
}

func TestFilter_ByKind(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Deployment", "api", "default", "ScalingUp", "scaled", 5*time.Minute),
		newEvent("Normal", "Service", "svc-1", "default", "Created", "created", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour, Kinds: []string{"Pod", "Deployment"}})
	if len(result) != 2 {
		t.Errorf("expected 2 events for Pod+Deployment, got %d", len(result))
	}
}

func TestFilter_ByName(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "api-server", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "worker", "default", "Scheduled", "scheduled", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour, Names: []string{"api-server"}})
	if len(result) != 1 {
		t.Errorf("expected 1 event for api-server, got %d", len(result))
	}
	if result[0].InvolvedObject.Name != "api-server" {
		t.Errorf("expected name api-server, got %s", result[0].InvolvedObject.Name)
	}
}

func TestFilter_ByType(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "app-2", "default", "Scheduled", "scheduled", 5*time.Minute),
		newEvent("Warning", "Pod", "app-3", "default", "Unhealthy", "probe failed", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour, Types: []string{"Warning"}})
	if len(result) != 2 {
		t.Errorf("expected 2 Warning events, got %d", len(result))
	}
}

func TestFilter_ByReason(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Warning", "Pod", "app-2", "default", "Unhealthy", "probe failed", 5*time.Minute),
		newEvent("Warning", "Pod", "app-3", "default", "BackOff", "back-off", 10*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour, Reasons: []string{"BackOff"}})
	if len(result) != 2 {
		t.Errorf("expected 2 BackOff events, got %d", len(result))
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour, Kinds: []string{"pod"}})
	if len(result) != 1 {
		t.Errorf("expected case-insensitive match, got %d", len(result))
	}
}

func TestFilter_Combined(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "api", "prod", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "api", "prod", "Scheduled", "scheduled", 5*time.Minute),
		newEvent("Warning", "Deployment", "api", "prod", "ScalingUp", "scaled", 5*time.Minute),
		newEvent("Warning", "Pod", "worker", "prod", "BackOff", "back-off", 5*time.Minute),
	}

	result := Filter(events, FilterOptions{
		Since: time.Hour,
		Kinds: []string{"Pod"},
		Names: []string{"api"},
		Types: []string{"Warning"},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 event (Warning+Pod+api), got %d", len(result))
	}
}

func TestFilter_Empty(t *testing.T) {
	result := Filter(nil, FilterOptions{Since: time.Hour})
	if len(result) != 0 {
		t.Errorf("expected 0 events for nil input, got %d", len(result))
	}
}

func TestFilter_SortedByTime(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "old", "default", "BackOff", "old", 30*time.Minute),
		newEvent("Warning", "Pod", "new", "default", "BackOff", "new", 1*time.Minute),
		newEvent("Warning", "Pod", "mid", "default", "BackOff", "mid", 15*time.Minute),
	}

	result := Filter(events, FilterOptions{Since: time.Hour})
	if len(result) != 3 {
		t.Fatalf("expected 3 events, got %d", len(result))
	}
	if result[0].InvolvedObject.Name != "new" {
		t.Errorf("expected newest first, got %s", result[0].InvolvedObject.Name)
	}
	if result[2].InvolvedObject.Name != "old" {
		t.Errorf("expected oldest last, got %s", result[2].InvolvedObject.Name)
	}
}

func TestGroupByResource(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off 1", 5*time.Minute),
		newEvent("Normal", "Pod", "app-1", "default", "Scheduled", "scheduled", 10*time.Minute),
		newEvent("Warning", "Pod", "app-2", "default", "Unhealthy", "probe failed", 3*time.Minute),
		newEvent("Normal", "Deployment", "api", "default", "ScalingUp", "scaled", 8*time.Minute),
	}

	groups := GroupByResource(events)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}

	// First group should have warnings (app-1 or app-2)
	if !hasWarningInGroup(groups[0]) {
		t.Error("expected first group to have warnings")
	}
}

func TestGroupByResource_WarningsFirst(t *testing.T) {
	events := []Event{
		newEvent("Normal", "Deployment", "api", "default", "ScalingUp", "scaled", 1*time.Minute),
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
	}

	groups := GroupByResource(events)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key.Kind != "Pod" {
		t.Errorf("expected warning group (Pod) first, got %s", groups[0].Key.Kind)
	}
}

func TestGroupByResource_Empty(t *testing.T) {
	groups := GroupByResource(nil)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for nil input, got %d", len(groups))
	}
}

func TestGroupByResource_SameResource(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app", "default", "BackOff", "msg1", 5*time.Minute),
		newEvent("Warning", "Pod", "app", "default", "Unhealthy", "msg2", 3*time.Minute),
		newEvent("Normal", "Pod", "app", "default", "Scheduled", "msg3", 10*time.Minute),
	}

	groups := GroupByResource(events)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Events) != 3 {
		t.Errorf("expected 3 events in group, got %d", len(groups[0].Events))
	}
}

func TestContainsCI(t *testing.T) {
	tests := []struct {
		list   []string
		target string
		want   bool
	}{
		{[]string{"Pod", "Deployment"}, "pod", true},
		{[]string{"Pod", "Deployment"}, "Pod", true},
		{[]string{"Pod"}, "Service", false},
		{nil, "Pod", false},
		{[]string{}, "Pod", false},
	}

	for _, tt := range tests {
		got := containsCI(tt.list, tt.target)
		if got != tt.want {
			t.Errorf("containsCI(%v, %q) = %v, want %v", tt.list, tt.target, got, tt.want)
		}
	}
}

func TestGroupEvents_Resource(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "app-1", "default", "Scheduled", "scheduled", 10*time.Minute),
		newEvent("Normal", "Deployment", "api", "prod", "ScalingUp", "scaled", 3*time.Minute),
	}

	groups := GroupEvents(events, GroupResource)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// Warning group should be first
	if !hasWarningInGroup(groups[0]) {
		t.Error("expected first group to have warnings")
	}
}

func TestGroupEvents_Namespace(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "prod", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "app-2", "prod", "Scheduled", "scheduled", 3*time.Minute),
		newEvent("Normal", "Deployment", "api", "staging", "ScalingUp", "scaled", 8*time.Minute),
	}

	groups := GroupEvents(events, GroupNamespace)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// prod has warnings, should be first
	if groups[0].Key.Label != "prod" {
		t.Errorf("expected first group label=prod, got %s", groups[0].Key.Label)
	}
}

func TestGroupEvents_Kind(t *testing.T) {
	events := []Event{
		newEvent("Normal", "Deployment", "api", "default", "ScalingUp", "scaled", 3*time.Minute),
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "app-2", "default", "Scheduled", "scheduled", 8*time.Minute),
	}

	groups := GroupEvents(events, GroupKind)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// Pod group has warnings, should be first
	if groups[0].Key.Label != "Pod" {
		t.Errorf("expected first group label=Pod, got %s", groups[0].Key.Label)
	}
}

func TestGroupEvents_Reason(t *testing.T) {
	events := []Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Warning", "Pod", "app-2", "default", "BackOff", "back-off again", 3*time.Minute),
		newEvent("Normal", "Pod", "app-3", "default", "Scheduled", "scheduled", 8*time.Minute),
	}

	groups := GroupEvents(events, GroupReason)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// BackOff has warnings, should be first
	if groups[0].Key.Label != "BackOff" {
		t.Errorf("expected first group label=BackOff, got %s", groups[0].Key.Label)
	}
	if len(groups[0].Events) != 2 {
		t.Errorf("expected 2 events in BackOff group, got %d", len(groups[0].Events))
	}
}

func TestGroupEvents_ClusterScoped(t *testing.T) {
	events := []Event{
		newEvent("Normal", "Node", "node-1", "", "NodeReady", "node is ready", 5*time.Minute),
	}

	groups := GroupEvents(events, GroupNamespace)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Key.Label != "(cluster-scoped)" {
		t.Errorf("expected label=(cluster-scoped), got %s", groups[0].Key.Label)
	}
}

func TestValidGroupBy(t *testing.T) {
	valid := []string{"", "resource", "namespace", "kind", "reason"}
	for _, v := range valid {
		if !ValidGroupBy(v) {
			t.Errorf("ValidGroupBy(%q) = false, want true", v)
		}
	}

	invalid := []string{"invalid", "node", "type"}
	for _, v := range invalid {
		if ValidGroupBy(v) {
			t.Errorf("ValidGroupBy(%q) = true, want false", v)
		}
	}
}

func hasWarningInGroup(g ResourceGroup) bool {
	for _, e := range g.Events {
		if e.Type == "Warning" {
			return true
		}
	}
	return false
}
