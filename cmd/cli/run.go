package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/somaz94/kube-events/internal/client"
	"github.com/somaz94/kube-events/internal/event"
	"github.com/somaz94/kube-events/internal/report"
	"github.com/spf13/cobra"
)

type eventFlags struct {
	kubeconfig    string
	kubeContext   string
	namespaces    []string
	kinds         []string
	names         []string
	types         []string
	reasons       []string
	since         string
	output        string
	summaryOnly   bool
	allNamespaces bool
	watch         bool
}

func extractFlags(cmd *cobra.Command) (eventFlags, error) {
	f := eventFlags{}
	var err error

	f.kubeconfig, err = cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return f, err
	}
	f.kubeContext, err = cmd.Flags().GetString("context")
	if err != nil {
		return f, err
	}
	f.namespaces, err = cmd.Flags().GetStringSlice("namespace")
	if err != nil {
		return f, err
	}
	f.kinds, err = cmd.Flags().GetStringSlice("kind")
	if err != nil {
		return f, err
	}
	f.names, err = cmd.Flags().GetStringSlice("name")
	if err != nil {
		return f, err
	}
	f.types, err = cmd.Flags().GetStringSlice("type")
	if err != nil {
		return f, err
	}
	f.reasons, err = cmd.Flags().GetStringSlice("reason")
	if err != nil {
		return f, err
	}
	f.since, err = cmd.Flags().GetString("since")
	if err != nil {
		return f, err
	}
	f.output, err = cmd.Flags().GetString("output")
	if err != nil {
		return f, err
	}
	f.summaryOnly, err = cmd.Flags().GetBool("summary-only")
	if err != nil {
		return f, err
	}
	f.allNamespaces, err = cmd.Flags().GetBool("all-namespaces")
	if err != nil {
		return f, err
	}
	f.watch, err = cmd.Flags().GetBool("watch")
	if err != nil {
		return f, err
	}

	return f, nil
}

func parseSince(s string) (time.Duration, error) {
	if s == "" {
		return time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", s, err)
	}
	return d, nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	f, err := extractFlags(cmd)
	if err != nil {
		return err
	}

	if f.watch {
		return runWatch(f)
	}

	since, err := parseSince(f.since)
	if err != nil {
		return err
	}

	c, err := client.New(f.kubeconfig, f.kubeContext)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ctx := context.Background()

	namespaces := f.namespaces
	if f.allNamespaces || len(namespaces) == 0 {
		namespaces = []string{""}
	}

	var allEvents []event.Event
	for _, ns := range namespaces {
		events, err := c.ListEvents(ctx, ns)
		if err != nil {
			return fmt.Errorf("failed to list events in namespace %q: %w", ns, err)
		}
		allEvents = append(allEvents, events...)
	}

	// Apply filters
	filtered := event.Filter(allEvents, event.FilterOptions{
		Since:   since,
		Kinds:   f.kinds,
		Names:   f.names,
		Types:   toUpper(f.types),
		Reasons: f.reasons,
	})

	// Group by resource
	groups := event.GroupByResource(filtered)

	// Build summary
	summary := report.NewSummary(groups, filtered)

	// Output
	w := os.Stdout
	switch f.output {
	case "json":
		return summary.PrintJSON(w)
	case "markdown":
		return summary.PrintMarkdown(w)
	case "table":
		return summary.PrintTable(w)
	case "plain":
		return summary.PrintPlain(w, f.summaryOnly)
	default:
		return summary.PrintColor(w, f.summaryOnly)
	}
}

func toUpper(ss []string) []string {
	result := make([]string, len(ss))
	for i, s := range ss {
		result[i] = strings.ToUpper(s[:1]) + s[1:]
	}
	return result
}
