package cli

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/somaz94/kube-events/internal/event"
	"github.com/somaz94/kube-events/internal/report"
	"github.com/spf13/cobra"
)

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

func TestRootCmd_FlagDefaults(t *testing.T) {
	// Flags() only merges PersistentFlags at execution time, so read them directly.
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

func TestRunEvents_TypeFilterIgnoresCase(t *testing.T) {
	now := time.Now()
	lister := &fakeLister{events: map[string][]event.Event{"": {
		{Type: "Warning", Reason: "BackOff", LastSeen: now,
			InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "a"}},
		{Type: "Normal", Reason: "Pulled", LastSeen: now,
			InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "b"}},
	}}}

	w, read := captureFile(t)
	f := eventFlags{output: "json", since: "1h", types: []string{"warning"}}
	if err := runEvents(lister, f, w); err != nil {
		t.Fatalf("runEvents error: %v", err)
	}
	out := read()
	if !strings.Contains(out, "BackOff") || strings.Contains(out, "Pulled") {
		t.Errorf("expected only the Warning event for --type warning, got:\n%s", out)
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

func TestPrintWatchEvent_JSON(t *testing.T) {
	e := event.Event{
		Type: "Warning", Reason: "BackOff", Message: "Back-off restarting", Count: 3,
		LastSeen: time.Now(), Age: 30 * time.Second,
		InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "app-1", Namespace: "default"},
	}

	w, read := captureFile(t)
	printWatchEvent(w, e, "json")

	var got struct {
		Groups []struct {
			Name   string `json:"name"`
			Events []struct {
				Reason string `json:"reason"`
			} `json:"events"`
		} `json:"groups"`
	}
	if err := json.Unmarshal([]byte(read()), &got); err != nil {
		t.Fatalf("output is not one JSON document: %v", err)
	}
	if len(got.Groups) != 1 || got.Groups[0].Name != "app-1" || got.Groups[0].Events[0].Reason != "BackOff" {
		t.Errorf("unexpected JSON: %+v", got)
	}
}

func TestPrintWatchEvent_ClusterScopedHasNoNamespaceTag(t *testing.T) {
	e := event.Event{
		Type: "Normal", Reason: "Pulled", Message: "Pulled image",
		LastSeen: time.Now(), Age: 5 * time.Minute,
		InvolvedObject: event.InvolvedObject{Kind: "Node", Name: "node-1"},
	}

	for _, format := range []string{"color", "plain"} {
		w, read := captureFile(t)
		printWatchEvent(w, e, format)
		out := read()
		if !strings.Contains(out, "Node/node-1") || strings.Contains(out, report.ColorCyan) || strings.Contains(out, " [") {
			t.Errorf("%s: expected no namespace tag, got %q", format, out)
		}
	}
}

func TestPrintWatchEvent_ColorOnlyForColorFormat(t *testing.T) {
	e := event.Event{
		Type: "Warning", Reason: "BackOff", Message: "Back-off restarting", Count: 1,
		LastSeen: time.Now(), Age: 30 * time.Second,
		InvolvedObject: event.InvolvedObject{Kind: "Pod", Name: "app-1", Namespace: "default"},
	}

	for _, format := range []string{"color", "plain", "markdown", "table"} {
		f, err := os.CreateTemp(t.TempDir(), "watch-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		printWatchEvent(f, e, format)
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		out := string(raw)

		wantColor := format == "color"
		if got := strings.Contains(out, "\033["); got != wantColor {
			t.Errorf("%s: ANSI codes present = %v, want %v: %q", format, got, wantColor, out)
		}
		if !strings.Contains(out, "BackOff") || !strings.Contains(out, "app-1") || !strings.Contains(out, "[default]") {
			t.Errorf("%s: expected reason, object and namespace in %q", format, out)
		}
	}
}

func TestExtractFlags_WithArgs(t *testing.T) {
	cmd := newFlagCmd()

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
	cmd := &cobra.Command{Use: "test"}
	_, err := extractFlags(cmd)
	if err == nil {
		t.Error("expected error for missing flags")
	}
}
