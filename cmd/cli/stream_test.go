package cli

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/somaz94/kube-events/internal/event"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// captureFile returns a temp file to stream into, plus a reader for its contents.
func captureFile(t *testing.T) (*os.File, func() string) {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "kube-events-stream-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f, func() string {
		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatalf("read captured output: %v", err)
		}
		return string(b)
	}
}

// k8sEvent builds a core/v1 Event recent enough to survive a 1h Since filter.
func k8sEvent(kind, name, namespace, reason, evType string) *corev1.Event {
	now := metav1.NewTime(time.Now())
	return &corev1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: name + ".event", Namespace: namespace},
		InvolvedObject: corev1.ObjectReference{Kind: kind, Name: name, Namespace: namespace},
		Reason:         reason,
		Message:        reason + " happened",
		Type:           evType,
		Count:          1,
		LastTimestamp:  now,
		FirstTimestamp: now,
	}
}

func defaultOpts() event.FilterOptions {
	return event.FilterOptions{Since: time.Hour}
}

func TestStreamEvents_PrintsAddedAndModified(t *testing.T) {
	w, read := captureFile(t)
	fake := watch.NewFake()

	go func() {
		fake.Add(k8sEvent("Pod", "added-pod", "default", "Scheduled", "Normal"))
		fake.Modify(k8sEvent("Pod", "modified-pod", "default", "BackOff", "Warning"))
		fake.Stop() // closes the channel, ending the loop
	}()

	if err := streamEvents(context.Background(), fake.ResultChan(), defaultOpts(), "plain", w); err != nil {
		t.Fatalf("streamEvents() error = %v, want nil", err)
	}

	out := read()
	for _, want := range []string{"added-pod", "modified-pod"} {
		if !strings.Contains(out, want) {
			t.Errorf("output = %q, want it to contain %q", out, want)
		}
	}
}

// Deleted and Bookmark events carry no new information for the stream and must
// not be printed.
func TestStreamEvents_IgnoresOtherWatchTypes(t *testing.T) {
	w, read := captureFile(t)
	fake := watch.NewFake()

	go func() {
		fake.Delete(k8sEvent("Pod", "deleted-pod", "default", "Killing", "Normal"))
		fake.Action(watch.Bookmark, k8sEvent("Pod", "bookmark-pod", "default", "Sync", "Normal"))
		fake.Add(k8sEvent("Pod", "kept-pod", "default", "Scheduled", "Normal"))
		fake.Stop()
	}()

	if err := streamEvents(context.Background(), fake.ResultChan(), defaultOpts(), "plain", w); err != nil {
		t.Fatalf("streamEvents() error = %v, want nil", err)
	}

	out := read()
	if !strings.Contains(out, "kept-pod") {
		t.Errorf("output = %q, want the Added event to be printed", out)
	}
	for _, unwanted := range []string{"deleted-pod", "bookmark-pod"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("output = %q, want it NOT to contain %q", out, unwanted)
		}
	}
}

// An object that is not a core/v1 Event must be skipped rather than panic.
func TestStreamEvents_SkipsNonEventObjects(t *testing.T) {
	w, read := captureFile(t)
	fake := watch.NewFake()

	go func() {
		fake.Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "not-an-event"}})
		fake.Add(k8sEvent("Pod", "real-event", "default", "Scheduled", "Normal"))
		fake.Stop()
	}()

	if err := streamEvents(context.Background(), fake.ResultChan(), defaultOpts(), "plain", w); err != nil {
		t.Fatalf("streamEvents() error = %v, want nil", err)
	}

	out := read()
	if strings.Contains(out, "not-an-event") {
		t.Errorf("output = %q, want the non-Event object skipped", out)
	}
	if !strings.Contains(out, "real-event") {
		t.Errorf("output = %q, want the Event still printed", out)
	}
}

func TestStreamEvents_AppliesFilters(t *testing.T) {
	tests := []struct {
		name    string
		opts    event.FilterOptions
		want    []string
		notWant []string
	}{
		{
			name:    "by kind",
			opts:    event.FilterOptions{Since: time.Hour, Kinds: []string{"Pod"}},
			want:    []string{"a-pod"},
			notWant: []string{"a-node"},
		},
		{
			name:    "by reason",
			opts:    event.FilterOptions{Since: time.Hour, Reasons: []string{"BackOff"}},
			want:    []string{"a-node"},
			notWant: []string{"a-pod"},
		},
		{
			name:    "by type",
			opts:    event.FilterOptions{Since: time.Hour, Types: []string{"Warning"}},
			want:    []string{"a-node"},
			notWant: []string{"a-pod"},
		},
		{
			name:    "by name",
			opts:    event.FilterOptions{Since: time.Hour, Names: []string{"a-pod"}},
			want:    []string{"a-pod"},
			notWant: []string{"a-node"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, read := captureFile(t)
			fake := watch.NewFake()

			go func() {
				fake.Add(k8sEvent("Pod", "a-pod", "default", "Scheduled", "Normal"))
				fake.Add(k8sEvent("Node", "a-node", "", "BackOff", "Warning"))
				fake.Stop()
			}()

			if err := streamEvents(context.Background(), fake.ResultChan(), tt.opts, "plain", w); err != nil {
				t.Fatalf("streamEvents() error = %v, want nil", err)
			}

			out := read()
			for _, s := range tt.want {
				if !strings.Contains(out, s) {
					t.Errorf("output = %q, want it to contain %q", out, s)
				}
			}
			for _, s := range tt.notWant {
				if strings.Contains(out, s) {
					t.Errorf("output = %q, want it NOT to contain %q", out, s)
				}
			}
		})
	}
}

func TestStreamEvents_JSONOutput(t *testing.T) {
	w, read := captureFile(t)
	fake := watch.NewFake()

	go func() {
		fake.Add(k8sEvent("Pod", "json-pod", "default", "Scheduled", "Normal"))
		fake.Stop()
	}()

	if err := streamEvents(context.Background(), fake.ResultChan(), defaultOpts(), "json", w); err != nil {
		t.Fatalf("streamEvents() error = %v, want nil", err)
	}

	out := read()
	if !strings.Contains(out, "json-pod") || !strings.Contains(out, "{") {
		t.Errorf("output = %q, want JSON mentioning the event", out)
	}
}

// A cancelled context ends the loop even while the stream stays open.
func TestStreamEvents_StopsOnContextCancel(t *testing.T) {
	w, _ := captureFile(t)
	fake := watch.NewFake()
	defer fake.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() { done <- streamEvents(ctx, fake.ResultChan(), defaultOpts(), "plain", w) }()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("streamEvents() error = %v, want nil on cancellation", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("streamEvents() did not return after the context was cancelled")
	}
}

// A closed stream ends the loop.
func TestStreamEvents_StopsOnClosedChannel(t *testing.T) {
	w, _ := captureFile(t)
	ch := make(chan watch.Event)
	close(ch)

	done := make(chan error, 1)
	go func() { done <- streamEvents(context.Background(), ch, defaultOpts(), "plain", w) }()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("streamEvents() error = %v, want nil on a closed stream", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("streamEvents() did not return after the stream closed")
	}
}
