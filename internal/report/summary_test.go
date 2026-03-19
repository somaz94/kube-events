package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/somaz94/kube-events/internal/event"
)

func newEvent(typ, kind, name, ns, reason, msg string, age time.Duration) event.Event {
	return event.Event{
		Type:    typ,
		Reason:  reason,
		Message: msg,
		Count:   1,
		LastSeen:  time.Now().Add(-age),
		FirstSeen: time.Now().Add(-age),
		Age:       age,
		InvolvedObject: event.InvolvedObject{
			Kind:      kind,
			Name:      name,
			Namespace: ns,
		},
	}
}

func sampleGroups() ([]event.ResourceGroup, []event.Event) {
	events := []event.Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "Back-off restarting failed container", 5*time.Minute),
		newEvent("Warning", "Pod", "app-1", "default", "Unhealthy", "Readiness probe failed", 8*time.Minute),
		newEvent("Normal", "Deployment", "api", "prod", "ScalingUp", "Scaled up replica set", 2*time.Minute),
	}

	groups := []event.ResourceGroup{
		{
			Key:    event.ResourceKey{Kind: "Pod", Name: "app-1", Namespace: "default"},
			Events: events[:2],
		},
		{
			Key:    event.ResourceKey{Kind: "Deployment", Name: "api", Namespace: "prod"},
			Events: events[2:],
		},
	}

	return groups, events
}

func TestNewSummary(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	if s.TotalEvents != 3 {
		t.Errorf("expected TotalEvents=3, got %d", s.TotalEvents)
	}
	if s.WarningCount != 2 {
		t.Errorf("expected WarningCount=2, got %d", s.WarningCount)
	}
	if s.NormalCount != 1 {
		t.Errorf("expected NormalCount=1, got %d", s.NormalCount)
	}
	if s.Resources != 2 {
		t.Errorf("expected Resources=2, got %d", s.Resources)
	}
}

func TestNewSummary_Empty(t *testing.T) {
	s := NewSummary(nil, nil, "")

	if s.TotalEvents != 0 {
		t.Errorf("expected TotalEvents=0, got %d", s.TotalEvents)
	}
	if s.Resources != 0 {
		t.Errorf("expected Resources=0, got %d", s.Resources)
	}
}

func TestPrintColor(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintColor(&buf, false)
	if err != nil {
		t.Fatalf("PrintColor error: %v", err)
	}

	out := buf.String()

	// Should contain resource headers
	if !strings.Contains(out, "Pod/app-1") {
		t.Error("expected Pod/app-1 in output")
	}
	if !strings.Contains(out, "Deployment/api") {
		t.Error("expected Deployment/api in output")
	}

	// Should contain event reasons
	if !strings.Contains(out, "BackOff") {
		t.Error("expected BackOff reason in output")
	}

	// Should contain summary
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line in output")
	}

	// Should contain ANSI color codes
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI color codes in output")
	}
}

func TestPrintColor_SummaryOnly(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintColor(&buf, true)
	if err != nil {
		t.Fatalf("PrintColor summaryOnly error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line")
	}
	if strings.Contains(out, "BackOff") {
		t.Error("summaryOnly should not contain event details")
	}
}

func TestPrintColor_NoEvents(t *testing.T) {
	s := NewSummary(nil, nil, "")

	var buf bytes.Buffer
	err := s.PrintColor(&buf, false)
	if err != nil {
		t.Fatalf("PrintColor error: %v", err)
	}

	if !strings.Contains(buf.String(), "No events found") {
		t.Error("expected 'No events found' message")
	}
}

func TestPrintPlain(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintPlain(&buf, false)
	if err != nil {
		t.Fatalf("PrintPlain error: %v", err)
	}

	out := buf.String()

	// Should NOT contain ANSI codes
	if strings.Contains(out, "\033[") {
		t.Error("plain output should not contain ANSI color codes")
	}

	if !strings.Contains(out, "Pod/app-1") {
		t.Error("expected Pod/app-1 in plain output")
	}
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary in plain output")
	}
}

func TestPrintJSON(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintJSON(&buf)
	if err != nil {
		t.Fatalf("PrintJSON error: %v", err)
	}

	// Should be valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	// Check summary fields
	summary, ok := parsed["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("expected summary field in JSON")
	}
	if summary["totalEvents"].(float64) != 3 {
		t.Errorf("expected totalEvents=3, got %v", summary["totalEvents"])
	}
	if summary["warningCount"].(float64) != 2 {
		t.Errorf("expected warningCount=2, got %v", summary["warningCount"])
	}

	// Check groups
	groupsJSON, ok := parsed["groups"].([]interface{})
	if !ok {
		t.Fatal("expected groups field in JSON")
	}
	if len(groupsJSON) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groupsJSON))
	}
}

func TestPrintMarkdown(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintMarkdown(&buf)
	if err != nil {
		t.Fatalf("PrintMarkdown error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "## Kubernetes Events Summary") {
		t.Error("expected markdown header")
	}
	if !strings.Contains(out, "| Type |") {
		t.Error("expected markdown table header")
	}
	if !strings.Contains(out, "|---") {
		t.Error("expected markdown table separator")
	}
	if !strings.Contains(out, "Warning") {
		t.Error("expected Warning in table")
	}
}

func TestPrintTable(t *testing.T) {
	groups, events := sampleGroups()
	s := NewSummary(groups, events, "resource")

	var buf bytes.Buffer
	err := s.PrintTable(&buf)
	if err != nil {
		t.Fatalf("PrintTable error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "TYPE") {
		t.Error("expected TYPE column header")
	}
	if !strings.Contains(out, "GROUP") {
		t.Error("expected GROUP column header")
	}
	if !strings.Contains(out, "REASON") {
		t.Error("expected REASON column header")
	}
	if !strings.Contains(out, "Total:") {
		t.Error("expected Total line")
	}
}

func TestPrintMarkdown_NoEvents(t *testing.T) {
	s := NewSummary(nil, nil, "")

	var buf bytes.Buffer
	err := s.PrintMarkdown(&buf)
	if err != nil {
		t.Fatalf("PrintMarkdown error: %v", err)
	}

	if !strings.Contains(buf.String(), "No events found") {
		t.Error("expected 'No events found' in empty markdown")
	}
}

func TestPrintColor_GroupByNamespace(t *testing.T) {
	events := []event.Event{
		newEvent("Warning", "Pod", "app-1", "prod", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Pod", "app-2", "staging", "Scheduled", "scheduled", 3*time.Minute),
	}
	groups := []event.ResourceGroup{
		{Key: event.ResourceKey{Label: "prod"}, Events: events[:1]},
		{Key: event.ResourceKey{Label: "staging"}, Events: events[1:]},
	}
	s := NewSummary(groups, events, "namespace")

	var buf bytes.Buffer
	if err := s.PrintColor(&buf, false); err != nil {
		t.Fatalf("PrintColor error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Error("expected 'prod' in output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected 'staging' in output")
	}
}

func TestPrintPlain_GroupByKind(t *testing.T) {
	events := []event.Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
		newEvent("Normal", "Deployment", "api", "default", "ScalingUp", "scaled", 3*time.Minute),
	}
	groups := []event.ResourceGroup{
		{Key: event.ResourceKey{Label: "Pod"}, Events: events[:1]},
		{Key: event.ResourceKey{Label: "Deployment"}, Events: events[1:]},
	}
	s := NewSummary(groups, events, "kind")

	var buf bytes.Buffer
	if err := s.PrintPlain(&buf, false); err != nil {
		t.Fatalf("PrintPlain error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Pod (1 events)") {
		t.Error("expected 'Pod (1 events)' in output")
	}
}

func TestPrintJSON_GroupByReason(t *testing.T) {
	events := []event.Event{
		newEvent("Warning", "Pod", "app-1", "default", "BackOff", "back-off", 5*time.Minute),
	}
	groups := []event.ResourceGroup{
		{Key: event.ResourceKey{Label: "BackOff"}, Events: events},
	}
	s := NewSummary(groups, events, "reason")

	var buf bytes.Buffer
	if err := s.PrintJSON(&buf); err != nil {
		t.Fatalf("PrintJSON error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	groupsJSON := parsed["groups"].([]interface{})
	first := groupsJSON[0].(map[string]interface{})
	if first["group"] != "BackOff" {
		t.Errorf("expected group=BackOff, got %v", first["group"])
	}
}

func TestPrintMarkdown_GroupByNamespace(t *testing.T) {
	events := []event.Event{
		newEvent("Warning", "Pod", "app-1", "prod", "BackOff", "back-off", 5*time.Minute),
	}
	groups := []event.ResourceGroup{
		{Key: event.ResourceKey{Label: "prod"}, Events: events},
	}
	s := NewSummary(groups, events, "namespace")

	var buf bytes.Buffer
	if err := s.PrintMarkdown(&buf); err != nil {
		t.Fatalf("PrintMarkdown error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "namespaces") {
		t.Error("expected 'namespaces' in markdown summary")
	}
}

func TestGroupNoun(t *testing.T) {
	tests := []struct {
		mode string
		want string
	}{
		{"namespace", "namespaces"},
		{"kind", "kinds"},
		{"reason", "reasons"},
		{"resource", "resources"},
		{"", "resources"},
	}
	for _, tt := range tests {
		s := NewSummary(nil, nil, tt.mode)
		got := s.groupNoun()
		if got != tt.want {
			t.Errorf("groupNoun(%q) = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{30 * time.Second, "30s"},
		{5 * time.Minute, "5m"},
		{2 * time.Hour, "2h"},
		{48 * time.Hour, "2d"},
	}

	for _, tt := range tests {
		got := event.FormatAge(tt.duration)
		if got != tt.want {
			t.Errorf("FormatAge(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exact", 5, "exact"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.max)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
		}
	}
}
