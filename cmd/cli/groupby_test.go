package cli

import (
	"context"
	"strings"
	"testing"

	"github.com/somaz94/kube-events/internal/event"
)

func TestValidateGroupBy(t *testing.T) {
	tests := []struct {
		name    string
		groupBy string
		wantErr bool
	}{
		{"empty falls back to the default", "", false},
		{"resource", "resource", false},
		{"namespace", "namespace", false},
		{"kind", "kind", false},
		{"reason", "reason", false},
		{"typo", "namesapce", true},
		{"unknown value", "pod", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroupBy(tt.groupBy)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateGroupBy(%q) error = %v, wantErr = %v", tt.groupBy, err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.groupBy) {
				t.Errorf("error = %q, want it to quote the offending value", err.Error())
			}
		})
	}
}

// countingLister records whether the API was reached, so the tests can assert
// that validation happens first.
type countingLister struct {
	calls int
}

func (c *countingLister) ListEvents(_ context.Context, _ string) ([]event.Event, error) {
	c.calls++
	return nil, nil
}

// An invalid --group-by must fail before any events are listed.
func TestRunEvents_InvalidGroupBySkipsTheAPI(t *testing.T) {
	lister := &countingLister{}
	w, _ := captureFile(t)

	err := runEvents(lister, eventFlags{output: "color", since: "1h", groupBy: "bogus"}, w)
	if err == nil {
		t.Fatal("runEvents() error = nil, want the invalid --group-by to fail")
	}
	if lister.calls != 0 {
		t.Errorf("lister called %d times, want 0 before validation passes", lister.calls)
	}
}

// The watch path rejects the same value, with the same message, even though it
// never acts on --group-by.
func TestRunWatch_RejectsInvalidGroupBy(t *testing.T) {
	err := runWatch(eventFlags{since: "5m", groupBy: "bogus", kubeconfig: missingKubeconfig(t)})
	if err == nil {
		t.Fatal("runWatch() error = nil, want the invalid --group-by to fail")
	}
	if !strings.Contains(err.Error(), "--group-by") {
		t.Errorf("error = %q, want it to report --group-by", err.Error())
	}
	// Reaching the kubeconfig would mean the flag was not validated first.
	if strings.Contains(err.Error(), "kubeconfig") {
		t.Errorf("error = %q, want --group-by validated before the kubeconfig is loaded", err.Error())
	}
}

// Both modes must reject a bad value identically, so the message cannot drift.
func TestGroupByRejectionIsIdenticalAcrossModes(t *testing.T) {
	const bad = "namesapce"

	w, _ := captureFile(t)
	listErr := runEvents(&countingLister{}, eventFlags{output: "color", since: "1h", groupBy: bad}, w)
	watchErr := runWatch(eventFlags{since: "5m", groupBy: bad, kubeconfig: missingKubeconfig(t)})

	if listErr == nil || watchErr == nil {
		t.Fatalf("both modes must fail: list=%v watch=%v", listErr, watchErr)
	}
	if listErr.Error() != watchErr.Error() {
		t.Errorf("messages differ:\n  list : %q\n  watch: %q", listErr.Error(), watchErr.Error())
	}
}

// A valid --group-by is accepted on the watch path even though it has no effect
// there, so existing command lines keep working.
func TestRunWatch_AcceptsValidGroupBy(t *testing.T) {
	for _, g := range []string{"", "resource", "namespace", "kind", "reason"} {
		err := runWatch(eventFlags{since: "5m", groupBy: g, kubeconfig: missingKubeconfig(t)})
		if err == nil {
			t.Fatalf("runWatch(groupBy=%q) error = nil, want it to proceed to the kubeconfig", g)
		}
		if !strings.Contains(err.Error(), "kubeconfig") {
			t.Errorf("runWatch(groupBy=%q) error = %q, want it past validation", g, err.Error())
		}
	}
}
