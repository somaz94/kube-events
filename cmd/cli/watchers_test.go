package cli

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/watch"
)

func TestResolveWatchNamespaces(t *testing.T) {
	tests := []struct {
		name  string
		flags eventFlags
		want  []string
	}{
		{"no namespace means cluster-wide", eventFlags{}, []string{""}},
		{"all-namespaces means cluster-wide", eventFlags{allNamespaces: true}, []string{""}},
		{"single namespace", eventFlags{namespaces: []string{"a"}}, []string{"a"}},
		{"every repeated namespace is kept", eventFlags{namespaces: []string{"a", "b", "c"}}, []string{"a", "b", "c"}},
		{
			"all-namespaces wins over an explicit list",
			eventFlags{namespaces: []string{"a", "b"}, allNamespaces: true},
			[]string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveWatchNamespaces(tt.flags)
			if len(got) != len(tt.want) {
				t.Fatalf("resolveWatchNamespaces() = %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("resolveWatchNamespaces()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// Every namespace on the command line must get its own watch — the regression
// this replaces silently watched only the first.
func TestStartWatchers_OpensOneWatchPerNamespace(t *testing.T) {
	var opened []string
	fakes := map[string]*watch.FakeWatcher{}
	start := func(_ context.Context, ns string) (watch.Interface, error) {
		opened = append(opened, ns)
		f := watch.NewFake()
		fakes[ns] = f
		return f, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	merged, stop, err := startWatchers(ctx, start, []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("startWatchers() error = %v, want nil", err)
	}
	defer stop()

	sort.Strings(opened)
	want := []string{"a", "b", "c"}
	if len(opened) != len(want) {
		t.Fatalf("opened watches on %v, want %v", opened, want)
	}
	for i := range want {
		if opened[i] != want[i] {
			t.Errorf("opened[%d] = %q, want %q", i, opened[i], want[i])
		}
	}

	// An event from any namespace must reach the merged stream.
	go func() {
		fakes["b"].Add(k8sEvent("Pod", "from-b", "b", "Scheduled", "Normal"))
	}()

	select {
	case ev := <-merged:
		if ev.Type != watch.Added {
			t.Errorf("merged event type = %v, want %v", ev.Type, watch.Added)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no event arrived on the merged stream")
	}
}

func TestStartWatchers_MergesEveryNamespace(t *testing.T) {
	fakes := map[string]*watch.FakeWatcher{}
	start := func(_ context.Context, ns string) (watch.Interface, error) {
		f := watch.NewFake()
		fakes[ns] = f
		return f, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	namespaces := []string{"a", "b"}
	merged, stop, err := startWatchers(ctx, start, namespaces)
	if err != nil {
		t.Fatalf("startWatchers() error = %v, want nil", err)
	}
	defer stop()

	go func() {
		fakes["a"].Add(k8sEvent("Pod", "pod-a", "a", "Scheduled", "Normal"))
		fakes["b"].Add(k8sEvent("Pod", "pod-b", "b", "Scheduled", "Normal"))
	}()

	seen := map[string]bool{}
	deadline := time.After(3 * time.Second)
	for len(seen) < 2 {
		select {
		case ev := <-merged:
			obj, ok := ev.Object.(interface{ GetNamespace() string })
			if !ok {
				t.Fatalf("unexpected object on the merged stream: %T", ev.Object)
			}
			seen[obj.GetNamespace()] = true
		case <-deadline:
			t.Fatalf("only saw namespaces %v, want both a and b", seen)
		}
	}
}

// A failure part-way through must stop the watchers already opened rather than
// leaking them.
func TestStartWatchers_StopsOpenedWatchersOnFailure(t *testing.T) {
	var opened []*watch.FakeWatcher
	wantErr := errors.New("forbidden")
	start := func(_ context.Context, ns string) (watch.Interface, error) {
		if ns == "bad" {
			return nil, wantErr
		}
		f := watch.NewFake()
		opened = append(opened, f)
		return f, nil
	}

	_, _, err := startWatchers(context.Background(), start, []string{"good", "bad"})
	if err == nil {
		t.Fatal("startWatchers() error = nil, want the failure to propagate")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}

	if len(opened) != 1 {
		t.Fatalf("opened %d watchers, want 1 before the failure", len(opened))
	}
	if !opened[0].IsStopped() {
		t.Error("the watcher opened before the failure is still running, want it stopped")
	}
}

// The merged stream closes once every underlying watcher has closed, which is
// what ends streamEvents.
func TestStartWatchers_MergedStreamClosesWithTheWatchers(t *testing.T) {
	fakes := map[string]*watch.FakeWatcher{}
	start := func(_ context.Context, ns string) (watch.Interface, error) {
		f := watch.NewFake()
		fakes[ns] = f
		return f, nil
	}

	merged, _, err := startWatchers(context.Background(), start, []string{"a", "b"})
	if err != nil {
		t.Fatalf("startWatchers() error = %v, want nil", err)
	}

	fakes["a"].Stop()
	fakes["b"].Stop()

	select {
	case _, ok := <-merged:
		if ok {
			t.Error("merged stream produced an event, want it closed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("merged stream did not close after every watcher stopped")
	}
}

func TestStartWatchers_ClusterWideOpensASingleWatch(t *testing.T) {
	var opened []string
	start := func(_ context.Context, ns string) (watch.Interface, error) {
		opened = append(opened, ns)
		return watch.NewFake(), nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, stop, err := startWatchers(ctx, start, resolveWatchNamespaces(eventFlags{allNamespaces: true}))
	if err != nil {
		t.Fatalf("startWatchers() error = %v, want nil", err)
	}
	defer stop()

	if len(opened) != 1 || opened[0] != "" {
		t.Errorf("opened %v, want a single cluster-wide watch", opened)
	}
}
