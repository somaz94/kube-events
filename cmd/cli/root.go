package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kube-events",
	Short: "A CLI tool to view and summarize Kubernetes events",
	Long: `kube-events provides a structured, color-coded view of Kubernetes events
grouped by resource with warning highlighting and summary statistics.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runRoot,
}

func init() {
	// Connection flags
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().String("context", "", "kubernetes context to use")

	// Filter flags
	rootCmd.PersistentFlags().StringSliceP("namespace", "n", nil, "filter by namespace (repeatable)")
	rootCmd.PersistentFlags().StringSliceP("kind", "k", nil, "filter by involved object kind (e.g., Pod, Deployment)")
	rootCmd.PersistentFlags().StringSliceP("name", "N", nil, "filter by involved object name")
	rootCmd.PersistentFlags().StringSliceP("type", "t", nil, "filter by event type (Normal, Warning)")
	rootCmd.PersistentFlags().StringSliceP("reason", "r", nil, "filter by event reason (e.g., BackOff, Unhealthy)")
	rootCmd.PersistentFlags().String("since", "1h", "show events newer than a relative duration (e.g., 5m, 1h, 24h)")

	// Output flags
	rootCmd.PersistentFlags().StringP("output", "o", "color", "output format: color, plain, json, markdown, table")
	rootCmd.PersistentFlags().BoolP("summary-only", "s", false, "show summary statistics only")
	rootCmd.PersistentFlags().Bool("all-namespaces", false, "show events from all namespaces")
	rootCmd.PersistentFlags().BoolP("watch", "w", false, "watch for new events in real-time")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
