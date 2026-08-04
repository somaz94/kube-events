package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newFlagCmd builds a command carrying the same flag set as the root command,
// so the wiring functions can be driven without touching the global rootCmd.
func newFlagCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("kubeconfig", "", "")
	cmd.Flags().String("context", "", "")
	cmd.Flags().StringSliceP("namespace", "n", nil, "")
	cmd.Flags().StringSliceP("kind", "k", nil, "")
	cmd.Flags().StringSliceP("name", "N", nil, "")
	cmd.Flags().StringSliceP("type", "t", nil, "")
	cmd.Flags().StringSliceP("reason", "r", nil, "")
	cmd.Flags().String("since", "1h", "")
	cmd.Flags().StringP("output", "o", "color", "")
	cmd.Flags().StringP("group-by", "g", "resource", "")
	cmd.Flags().BoolP("summary-only", "s", false, "")
	cmd.Flags().Bool("all-namespaces", false, "")
	cmd.Flags().BoolP("watch", "w", false, "")
	return cmd
}

// missingKubeconfig returns a path that is guaranteed not to exist, so loading
// it fails deterministically regardless of the developer's own kube config.
func missingKubeconfig(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "no-such-kubeconfig.yaml")
}

func TestRunRoot_ClientCreationFailure(t *testing.T) {
	cmd := newFlagCmd()
	if err := cmd.Flags().Set("kubeconfig", missingKubeconfig(t)); err != nil {
		t.Fatal(err)
	}

	err := runRoot(cmd, nil)
	if err == nil {
		t.Fatal("runRoot() error = nil, want the client creation to fail")
	}
	if !strings.Contains(err.Error(), "kubernetes client") {
		t.Errorf("error = %q, want it to name the failing step", err.Error())
	}
}

// --watch must reach runWatch, and an invalid --since must fail there before any
// connection is attempted.
func TestRunRoot_DispatchesToWatch(t *testing.T) {
	cmd := newFlagCmd()
	for flag, value := range map[string]string{
		"watch":      "true",
		"since":      "not-a-duration",
		"kubeconfig": missingKubeconfig(t),
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatal(err)
		}
	}

	err := runRoot(cmd, nil)
	if err == nil {
		t.Fatal("runRoot() error = nil, want the invalid --since to fail")
	}
	if !strings.Contains(err.Error(), "--since") {
		t.Errorf("error = %q, want the watch path to report --since", err.Error())
	}
	// Reaching the kubeconfig would mean the flags were not validated first.
	if strings.Contains(err.Error(), "kubeconfig") {
		t.Errorf("error = %q, want --since validated before the kubeconfig is loaded", err.Error())
	}
}

func TestRunRoot_FlagExtractionFailure(t *testing.T) {
	// A command missing the expected flags makes extractFlags fail.
	cmd := &cobra.Command{Use: "bare"}

	if err := runRoot(cmd, nil); err == nil {
		t.Fatal("runRoot() error = nil, want the flag extraction to fail")
	}
}

// The same ordering guarantee, asserted directly on runWatch.
func TestRunWatch_ValidatesSinceBeforeConnecting(t *testing.T) {
	err := runWatch(eventFlags{since: "nope", kubeconfig: missingKubeconfig(t)})
	if err == nil {
		t.Fatal("runWatch() error = nil, want the invalid --since to fail")
	}
	if !strings.Contains(err.Error(), "--since") {
		t.Errorf("error = %q, want it to report --since", err.Error())
	}
}

func TestRunWatch_KubeconfigFailure(t *testing.T) {
	err := runWatch(eventFlags{since: "5m", kubeconfig: missingKubeconfig(t)})
	if err == nil {
		t.Fatal("runWatch() error = nil, want loading the kubeconfig to fail")
	}
	if !strings.Contains(err.Error(), "kubeconfig") {
		t.Errorf("error = %q, want it to name the kubeconfig", err.Error())
	}
}

func TestExecute_RunsTheRootCommand(t *testing.T) {
	original := rootCmd.Flags().Args()
	t.Cleanup(func() { rootCmd.SetArgs(original) })

	rootCmd.SetArgs([]string{"version"})
	if err := Execute(); err != nil {
		t.Errorf("Execute() error = %v, want the version subcommand to succeed", err)
	}
}
