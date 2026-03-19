package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/somaz94/kube-events/internal/event"
	"github.com/somaz94/kube-events/internal/report"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func runWatch(f eventFlags) error {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if f.kubeconfig != "" {
		rules.ExplicitPath = f.kubeconfig
	}

	overrides := &clientcmd.ConfigOverrides{}
	if f.kubeContext != "" {
		overrides.CurrentContext = f.kubeContext
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
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

	since, _ := parseSince(f.since)
	filterOpts := event.FilterOptions{
		Since:   since,
		Kinds:   f.kinds,
		Names:   f.names,
		Types:   toUpper(f.types),
		Reasons: f.reasons,
	}

	fmt.Fprintf(os.Stderr, "Watching events (press Ctrl+C to stop)...\n\n")

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-watcher.ResultChan():
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

			e := convertWatchEvent(*k8sEvent)

			// Apply filters
			filtered := event.Filter([]event.Event{e}, filterOpts)
			if len(filtered) == 0 {
				continue
			}

			printWatchEvent(os.Stdout, filtered[0], f.output)
		}
	}
}

func convertWatchEvent(e corev1.Event) event.Event {
	lastSeen := e.LastTimestamp.Time
	if lastSeen.IsZero() {
		lastSeen = e.EventTime.Time
	}
	if lastSeen.IsZero() {
		lastSeen = e.CreationTimestamp.Time
	}

	return event.Event{
		Type:    e.Type,
		Reason:  e.Reason,
		Message: e.Message,
		Count:   e.Count,
		LastSeen:  lastSeen,
		FirstSeen: e.FirstTimestamp.Time,
		Age:       time.Since(lastSeen),
		InvolvedObject: event.InvolvedObject{
			Kind:      e.InvolvedObject.Kind,
			Name:      e.InvolvedObject.Name,
			Namespace: e.InvolvedObject.Namespace,
		},
		Source: event.Source{
			Component: e.Source.Component,
			Host:      e.Source.Host,
		},
	}
}

func printWatchEvent(w *os.File, e event.Event, format string) {
	switch format {
	case "json":
		groups := []event.ResourceGroup{{
			Key:    event.ResourceKey{Kind: e.InvolvedObject.Kind, Name: e.InvolvedObject.Name, Namespace: e.InvolvedObject.Namespace},
			Events: []event.Event{e},
		}}
		s := report.NewSummary(groups, []event.Event{e})
		s.PrintJSON(w)
	default:
		typeColor := "\033[32m"
		icon := "  "
		if e.Type == "Warning" {
			typeColor = "\033[33m"
			icon = "! "
		}

		ns := ""
		if e.InvolvedObject.Namespace != "" {
			ns = fmt.Sprintf(" \033[36m[%s]\033[0m", e.InvolvedObject.Namespace)
		}

		age := formatWatchAge(e.Age)
		fmt.Fprintf(w, "%s%s%-18s\033[0m %s%-8s\033[0m \033[1m%s/%s\033[0m%s %s\n",
			typeColor, icon, e.Reason,
			"\033[90m", age,
			e.InvolvedObject.Kind, e.InvolvedObject.Name, ns,
			e.Message)
	}
}

func formatWatchAge(d time.Duration) string {
	sec := d.Seconds()
	switch {
	case sec < 60:
		return fmt.Sprintf("%ds", int(sec))
	case sec < 3600:
		return fmt.Sprintf("%dm", int(sec/60))
	default:
		return fmt.Sprintf("%dh", int(sec/3600))
	}
}

