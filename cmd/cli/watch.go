package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/somaz94/kube-events/internal/client"
	"github.com/somaz94/kube-events/internal/event"
	"github.com/somaz94/kube-events/internal/report"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func runWatch(f eventFlags) error {
	config, err := client.LoadKubeConfig(f.kubeconfig, f.kubeContext)
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace := ""
	if len(f.namespaces) > 0 && !f.allNamespaces {
		namespace = f.namespaces[0]
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

	watcher, err := cs.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch events: %w", err)
	}
	defer watcher.Stop()

	since, err := parseSince(f.since)
	if err != nil {
		return fmt.Errorf("invalid --since value: %w", err)
	}
	filterOpts := event.FilterOptions{
		Since:   since,
		Kinds:   f.kinds,
		Names:   f.names,
		Types:   toUpper(f.types),
		Reasons: f.reasons,
	}

	fmt.Fprintf(os.Stderr, "Watching events (press Ctrl+C to stop)...\n\n")

	return streamEvents(ctx, watcher.ResultChan(), filterOpts, f.output, os.Stdout)
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
