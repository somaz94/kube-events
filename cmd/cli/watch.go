package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/somaz94/kube-events/internal/client"
	"github.com/somaz94/kube-events/internal/event"
	"github.com/somaz94/kube-events/internal/report"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// watchFunc opens a watch stream for a single namespace. An empty namespace
// means cluster-wide. It is the seam that lets tests supply fake streams in
// place of a live API server.
type watchFunc func(ctx context.Context, namespace string) (watch.Interface, error)

// resolveWatchNamespaces returns the namespaces to open a watch on. A single
// empty entry means cluster-wide, matching how runEvents treats the same flags.
func resolveWatchNamespaces(f eventFlags) []string {
	if f.allNamespaces || len(f.namespaces) == 0 {
		return []string{""}
	}
	return f.namespaces
}

// startWatchers opens one watcher per namespace and merges their streams into a
// single channel.
//
// The Kubernetes watch API is scoped to one namespace per call, so a repeatable
// --namespace can only be honored by fanning in; watching cluster-wide and
// filtering afterwards would instead demand cluster-scoped RBAC the caller may
// not have. The returned stop function stops every watcher that was opened.
func startWatchers(ctx context.Context, start watchFunc, namespaces []string) (<-chan watch.Event, func(), error) {
	watchers := make([]watch.Interface, 0, len(namespaces))
	stop := func() {
		for _, w := range watchers {
			w.Stop()
		}
	}

	for _, ns := range namespaces {
		w, err := start(ctx, ns)
		if err != nil {
			// Stop the watchers already opened so a partial failure leaks none.
			stop()
			return nil, nil, fmt.Errorf("failed to watch events in namespace %q: %w", ns, err)
		}
		watchers = append(watchers, w)
	}

	merged := make(chan watch.Event)
	var wg sync.WaitGroup
	for _, w := range watchers {
		wg.Add(1)
		go func(w watch.Interface) {
			defer wg.Done()
			for ev := range w.ResultChan() {
				select {
				case merged <- ev:
				case <-ctx.Done():
					return
				}
			}
		}(w)
	}
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged, stop, nil
}

func runWatch(f eventFlags) error {
	// Validate the flags before opening any connection, so a bad --since fails
	// immediately instead of after a watch has been established.
	since, err := parseSince(f.since)
	if err != nil {
		return fmt.Errorf("invalid --since value: %w", err)
	}

	config, err := client.LoadKubeConfig(f.kubeconfig, f.kubeContext)
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	start := func(ctx context.Context, namespace string) (watch.Interface, error) {
		return cs.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{})
	}

	merged, stop, err := startWatchers(ctx, start, resolveWatchNamespaces(f))
	if err != nil {
		return err
	}
	defer stop()

	filterOpts := event.FilterOptions{
		Since:   since,
		Kinds:   f.kinds,
		Names:   f.names,
		Types:   toUpper(f.types),
		Reasons: f.reasons,
	}

	fmt.Fprintf(os.Stderr, "Watching events (press Ctrl+C to stop)...\n\n")

	return streamEvents(ctx, merged, filterOpts, f.output, os.Stdout)
}

// streamEvents consumes a watch stream, applies the filters and prints every
// surviving event. It returns when the context is cancelled or the stream
// closes.
//
// Taking the channel rather than the watcher keeps the loop independent of how
// the stream was opened, so tests can drive it from watch.NewFake().
func streamEvents(
	ctx context.Context,
	ch <-chan watch.Event,
	opts event.FilterOptions,
	output string,
	w *os.File,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			if ev.Type != watch.Added && ev.Type != watch.Modified {
				continue
			}

			k8sEvent, ok := ev.Object.(*corev1.Event)
			if !ok {
				continue
			}

			e := event.ConvertK8sEvent(*k8sEvent)

			// Apply filters
			filtered := event.Filter([]event.Event{e}, opts)
			if len(filtered) == 0 {
				continue
			}

			printWatchEvent(w, filtered[0], output)
		}
	}
}

func printWatchEvent(w *os.File, e event.Event, format string) {
	switch format {
	case "json":
		groups := []event.ResourceGroup{{
			Key:    event.ResourceKey{Kind: e.InvolvedObject.Kind, Name: e.InvolvedObject.Name, Namespace: e.InvolvedObject.Namespace},
			Events: []event.Event{e},
		}}
		s := report.NewSummary(groups, []event.Event{e}, "resource")
		s.PrintJSON(w)
	default:
		typeColor, icon := report.EventStyle(e.Type)

		ns := ""
		if e.InvolvedObject.Namespace != "" {
			ns = fmt.Sprintf(" %s[%s]%s", report.ColorCyan, e.InvolvedObject.Namespace, report.ColorReset)
		}

		age := event.FormatAge(e.Age)
		fmt.Fprintf(w, "%s%s%-18s%s %s%-8s%s %s%s/%s%s%s %s\n",
			typeColor, icon, e.Reason, report.ColorReset,
			report.ColorGray, age, report.ColorReset,
			report.ColorBold, e.InvolvedObject.Kind, e.InvolvedObject.Name, report.ColorReset, ns,
			e.Message)
	}
}
