package cli

import (
	"testing"
	"time"
)

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
