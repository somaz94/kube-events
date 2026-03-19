package cli

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/somaz94/kube-events/internal/event"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// fakeLister implements client.EventLister for testing.
type fakeLister struct {
	events map[string][]event.Event
	err    error
}

func (f *fakeLister) ListEvents(_ context.Context, namespace string) ([]event.Event, error) {
	if f.err != nil {
		return nil, f.err
	}
	if namespace == "" {
		var all []event.Event
		for _, evts := range f.events {
			all = append(all, evts...)
		}
		return all, nil
	}
	return f.events[namespace], nil
}

func TestParseSince(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
		err   bool
	}{
		{"5m", 5 * time.Minute, false},
		{"1h", time.Hour, false},
		{"24h", 24 * time.Hour, false},
		{"30s", 30 * time.Second, false},
		{"", time.Hour, false},
		{"invalid", 0, true},
		{"abc123", 0, true},
	}

	for _, tt := range tests {
		got, err := parseSince(tt.input)
		if tt.err && err == nil {
			t.Errorf("parseSince(%q) expected error, got nil", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("parseSince(%q) unexpected error: %v", tt.input, err)
		}
		if !tt.err && got != tt.want {
			t.Errorf("parseSince(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input []string
		want  []string
	}{
		{[]string{"warning"}, []string{"Warning"}},
		{[]string{"normal", "warning"}, []string{"Normal", "Warning"}},
		{[]string{"Warning"}, []string{"Warning"}},
	}

	for _, tt := range tests {
		got := toUpper(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("toUpper(%v) length = %d, want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("toUpper(%v)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestExtractFlags_Defaults(t *testing.T) {
	// PersistentFlags are merged into Flags() during command execution.
	// Verify defaults via PersistentFlags directly.
	pf := rootCmd.PersistentFlags()

	since, _ := pf.GetString("since")
	if since != "1h" {
		t.Errorf("expected default since=1h, got %s", since)
	}
	output, _ := pf.GetString("output")
	if output != "color" {
		t.Errorf("expected default output=color, got %s", output)
	}
	groupBy, _ := pf.GetString("group-by")
	if groupBy != "resource" {
		t.Errorf("expected default group-by=resource, got %s", groupBy)
	}
	summaryOnly, _ := pf.GetBool("summary-only")
	if summaryOnly {
		t.Error("expected default summaryOnly=false")
	}
	allNs, _ := pf.GetBool("all-namespaces")
	if allNs {
		t.Error("expected default allNamespaces=false")
	}
	watch, _ := pf.GetBool("watch")
	if watch {
		t.Error("expected default watch=false")
	}
}

func TestRootCommandFlags(t *testing.T) {
	flags := []struct {
		name     string
		short    string
		hasShort bool
	}{
		{"kubeconfig", "", false},
		{"context", "", false},
		{"namespace", "n", true},
		{"kind", "k", true},
		{"name", "N", true},
		{"type", "t", true},
		{"reason", "r", true},
		{"since", "", false},
		{"output", "o", true},
		{"summary-only", "s", true},
		{"all-namespaces", "", false},
		{"watch", "w", true},
		{"group-by", "g", true},
	}

	for _, f := range flags {
		flag := rootCmd.PersistentFlags().Lookup(f.name)
		if flag == nil {
			t.Errorf("flag --%s not found", f.name)
			continue
		}
		if f.hasShort && flag.Shorthand != f.short {
			t.Errorf("flag --%s shorthand = %q, want %q", f.name, flag.Shorthand, f.short)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("version subcommand not found")
	}
}

func TestRunEvents_ColorOutput(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"default": {
				{Type: "Warning", Reason: "BackOff", Message: "Back-off restarting", Count: 3,
					LastSeen: now.Add(-2 * time.Minute), FirstSeen: now.Add(-10 * time.Minute), Age: 2 * time.Minute,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "app-1", Namespace: "default"},
					Source:         event.Source{Component: "kubelet", Host: "node-1"}},
				{Type: "Normal", Reason: "Scheduled", Message: "Successfully assigned", Count: 1,
					LastSeen: now.Add(-5 * time.Minute), FirstSeen: now.Add(-5 * time.Minute), Age: 5 * time.Minute,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "app-1", Namespace: "default"}},
			},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "color", since: "1h"}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(color) error: %v", err)
	}
}

func TestRunEvents_AllFormats(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {
				{Type: "Warning", Reason: "Unhealthy", Message: "Readiness probe failed", Count: 2,
					LastSeen: now.Add(-1 * time.Minute), FirstSeen: now.Add(-5 * time.Minute), Age: 1 * time.Minute,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "web-1", Namespace: "prod"}},
			},
		},
	}

	formats := []string{"json", "markdown", "table", "plain", "color"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			f := eventFlags{output: format, since: "1h", allNamespaces: true}
			if err := runEvents(lister, f, tmpFile); err != nil {
				t.Fatalf("runEvents(%s) error: %v", format, err)
			}
		})
	}
}

func TestRunEvents_SummaryOnly(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {{Type: "Normal", Reason: "Pulled", Message: "Pulled image", Count: 1,
				LastSeen: now, FirstSeen: now, Age: 0,
				InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "x", Namespace: "default"}}},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "plain", since: "1h", summaryOnly: true}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(summary-only) error: %v", err)
	}
}

func TestRunEvents_NamespaceFilter(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"prod": {{Type: "Warning", Reason: "OOM", Message: "OOMKilled", Count: 1,
				LastSeen: now, FirstSeen: now, Age: 0,
				InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "api", Namespace: "prod"}}},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "json", since: "1h", namespaces: []string{"prod"}}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(namespace) error: %v", err)
	}
}

func TestRunEvents_WithFilters(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {
				{Type: "Warning", Reason: "BackOff", Message: "Back-off", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a", Namespace: "default"}},
				{Type: "Normal", Reason: "Pulled", Message: "Pulled", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Deployment", Name: "b", Namespace: "default"}},
			},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{
		output:  "table",
		since:   "1h",
		kinds:   []string{"Pod"},
		names:   []string{"a"},
		types:   []string{"warning"},
		reasons: []string{"BackOff"},
	}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(filters) error: %v", err)
	}
}

func TestRunEvents_InvalidSince(t *testing.T) {
	lister := &fakeLister{events: map[string][]event.Event{}}
	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "color", since: "invalid"}
	if err := runEvents(lister, f, tmpFile); err == nil {
		t.Error("expected error for invalid since, got nil")
	}
}

func TestRunEvents_ListerError(t *testing.T) {
	lister := &fakeLister{err: context.DeadlineExceeded}
	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "color", since: "1h"}
	if err := runEvents(lister, f, tmpFile); err == nil {
		t.Error("expected error from lister, got nil")
	}
}

func TestRunEvents_Empty(t *testing.T) {
	lister := &fakeLister{events: map[string][]event.Event{}}
	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "json", since: "1h"}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(empty) error: %v", err)
	}
}

func TestConvertWatchEvent(t *testing.T) {
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

	e := event.ConvertK8sEvent(k8sEvent)
	if e.Type != "Warning" {
		t.Errorf("expected Warning, got %s", e.Type)
	}
	if e.Reason != "BackOff" {
		t.Errorf("expected BackOff, got %s", e.Reason)
	}
	if e.InvolvedObject.Kind != "Pod" {
		t.Errorf("expected Pod, got %s", e.InvolvedObject.Kind)
	}
	if e.Source.Component != "kubelet" {
		t.Errorf("expected kubelet, got %s", e.Source.Component)
	}
}

func TestConvertWatchEvent_FallbackTimestamps(t *testing.T) {
	now := time.Now()

	// EventTime fallback
	e1 := event.ConvertK8sEvent(corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "e1"},
		EventTime:  metav1.MicroTime{Time: now.Add(-1 * time.Minute)},
	})
	if e1.LastSeen.IsZero() {
		t.Error("expected LastSeen from EventTime")
	}

	// CreationTimestamp fallback
	e2 := event.ConvertK8sEvent(corev1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "e2", CreationTimestamp: metav1.Time{Time: now}},
	})
	if e2.LastSeen.IsZero() {
		t.Error("expected LastSeen from CreationTimestamp")
	}
}

func TestPrintWatchEvent(t *testing.T) {
	now := time.Now()
	e := event.Event{
		Type: "Warning", Reason: "BackOff", Message: "Back-off restarting", Count: 3,
		LastSeen: now, Age: 30 * time.Second,
		InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "app-1", Namespace: "default"},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-watch-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Default format
	printWatchEvent(tmpFile, e, "color")

	// JSON format
	tmpFile2, _ := os.CreateTemp("", "kube-events-watch-*.txt")
	defer os.Remove(tmpFile2.Name())
	defer tmpFile2.Close()
	printWatchEvent(tmpFile2, e, "json")
}

func TestPrintWatchEvent_Normal(t *testing.T) {
	e := event.Event{
		Type: "Normal", Reason: "Pulled", Message: "Pulled image", Count: 1,
		LastSeen: time.Now(), Age: 5 * time.Minute,
		InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "web-1"},
	}

	tmpFile, _ := os.CreateTemp("", "kube-events-watch-*.txt")
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	printWatchEvent(tmpFile, e, "color")
}

func TestExtractFlags_WithArgs(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.PersistentFlags().String("kubeconfig", "", "")
	cmd.PersistentFlags().String("context", "", "")
	cmd.PersistentFlags().StringSliceP("namespace", "n", nil, "")
	cmd.PersistentFlags().StringSliceP("kind", "k", nil, "")
	cmd.PersistentFlags().StringSliceP("name", "N", nil, "")
	cmd.PersistentFlags().StringSliceP("type", "t", nil, "")
	cmd.PersistentFlags().StringSliceP("reason", "r", nil, "")
	cmd.PersistentFlags().String("since", "1h", "")
	cmd.PersistentFlags().StringP("output", "o", "color", "")
	cmd.PersistentFlags().BoolP("summary-only", "s", false, "")
	cmd.PersistentFlags().Bool("all-namespaces", false, "")
	cmd.PersistentFlags().BoolP("watch", "w", false, "")
	cmd.PersistentFlags().StringP("group-by", "g", "resource", "")

	cmd.SetArgs([]string{
		"--kubeconfig", "/tmp/kc",
		"--context", "my-ctx",
		"-n", "prod",
		"-k", "Pod",
		"-N", "web-1",
		"-t", "Warning",
		"-r", "BackOff",
		"--since", "5m",
		"-o", "json",
		"-g", "namespace",
		"-s",
		"--all-namespaces",
		"-w",
	})

	var captured eventFlags
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error
		captured, err = extractFlags(cmd)
		return err
	}

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute error: %v", err)
	}

	if captured.kubeconfig != "/tmp/kc" {
		t.Errorf("kubeconfig = %q, want /tmp/kc", captured.kubeconfig)
	}
	if captured.kubeContext != "my-ctx" {
		t.Errorf("kubeContext = %q, want my-ctx", captured.kubeContext)
	}
	if len(captured.namespaces) != 1 || captured.namespaces[0] != "prod" {
		t.Errorf("namespaces = %v, want [prod]", captured.namespaces)
	}
	if len(captured.kinds) != 1 || captured.kinds[0] != "Pod" {
		t.Errorf("kinds = %v, want [Pod]", captured.kinds)
	}
	if captured.output != "json" {
		t.Errorf("output = %q, want json", captured.output)
	}
	if captured.groupBy != "namespace" {
		t.Errorf("groupBy = %q, want namespace", captured.groupBy)
	}
	if !captured.summaryOnly {
		t.Error("expected summaryOnly=true")
	}
	if !captured.allNamespaces {
		t.Error("expected allNamespaces=true")
	}
	if !captured.watch {
		t.Error("expected watch=true")
	}
}

func TestRunEvents_GroupByNamespace(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {
				{Type: "Warning", Reason: "BackOff", Message: "back-off", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a", Namespace: "prod"}},
				{Type: "Normal", Reason: "Pulled", Message: "pulled", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "b", Namespace: "staging"}},
			},
		},
	}

	formats := []string{"color", "plain", "json", "markdown", "table"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			f := eventFlags{output: format, since: "1h", groupBy: "namespace"}
			if err := runEvents(lister, f, tmpFile); err != nil {
				t.Fatalf("runEvents(%s, group-by=namespace) error: %v", format, err)
			}
		})
	}
}

func TestRunEvents_GroupByKind(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {
				{Type: "Warning", Reason: "BackOff", Message: "back-off", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a", Namespace: "default"}},
				{Type: "Normal", Reason: "ScalingUp", Message: "scaled", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Deployment", Name: "b", Namespace: "default"}},
			},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "json", since: "1h", groupBy: "kind"}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(group-by=kind) error: %v", err)
	}
}

func TestRunEvents_GroupByReason(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {
				{Type: "Warning", Reason: "BackOff", Message: "back-off", Count: 1,
					LastSeen: now, FirstSeen: now, Age: 0,
					InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a", Namespace: "default"}},
			},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "table", since: "1h", groupBy: "reason"}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(group-by=reason) error: %v", err)
	}
}

func TestRunEvents_InvalidGroupBy(t *testing.T) {
	lister := &fakeLister{events: map[string][]event.Event{}}
	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "color", since: "1h", groupBy: "invalid"}
	if err := runEvents(lister, f, tmpFile); err == nil {
		t.Error("expected error for invalid group-by, got nil")
	}
}

func TestToUpper_EmptyString(t *testing.T) {
	result := toUpper([]string{"", "warning", ""})
	if result[0] != "" {
		t.Errorf("expected empty string, got %q", result[0])
	}
	if result[1] != "Warning" {
		t.Errorf("expected Warning, got %q", result[1])
	}
	if result[2] != "" {
		t.Errorf("expected empty string, got %q", result[2])
	}
}

func TestToUpper_Nil(t *testing.T) {
	result := toUpper(nil)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d", len(result))
	}
}

func TestRunEvents_MultipleNamespaces(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"ns1": {{Type: "Warning", Reason: "BackOff", Message: "msg1", Count: 1,
				LastSeen: now, FirstSeen: now, Age: 0,
				InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a", Namespace: "ns1"}}},
			"ns2": {{Type: "Normal", Reason: "Pulled", Message: "msg2", Count: 1,
				LastSeen: now, FirstSeen: now, Age: 0,
				InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "b", Namespace: "ns2"}}},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "json", since: "1h", namespaces: []string{"ns1", "ns2"}}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(multi-ns) error: %v", err)
	}
}

func TestRunEvents_ColorSummaryOnly(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{
		events: map[string][]event.Event{
			"": {{Type: "Normal", Reason: "Pulled", Message: "Pulled", Count: 1,
				LastSeen: now, FirstSeen: now, Age: 0,
				InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "x", Namespace: "default"}}},
		},
	}

	tmpFile, err := os.CreateTemp("", "kube-events-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	f := eventFlags{output: "color", since: "1h", summaryOnly: true}
	if err := runEvents(lister, f, tmpFile); err != nil {
		t.Fatalf("runEvents(color-summary) error: %v", err)
	}
}

func TestExtractFlags_MissingFlag(t *testing.T) {
	// A command with no flags registered should fail
	cmd := &cobra.Command{Use: "test"}
	_, err := extractFlags(cmd)
	if err == nil {
		t.Error("expected error for missing flags")
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "30s"},
		{5 * time.Minute, "5m"},
		{2 * time.Hour, "2h"},
		{48 * time.Hour, "2d"},
		{0, "0s"},
		{59 * time.Second, "59s"},
		{60 * time.Second, "1m"},
		{3600 * time.Second, "1h"},
	}

	for _, tt := range tests {
		got := event.FormatAge(tt.d)
		if got != tt.want {
			t.Errorf("FormatAge(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
