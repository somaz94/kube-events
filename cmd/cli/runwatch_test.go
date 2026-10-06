package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeAPIServer starts an API server stand-in and returns a kubeconfig that points at it.
func fakeAPIServer(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	kubeconfig := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- name: fake
  cluster:
    server: %s
contexts:
- name: fake
  context:
    cluster: fake
    user: fake
current-context: fake
users:
- name: fake
  user: {}
`, srv.URL)
	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, []byte(kubeconfig), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// Each watch sends one event for its namespace and then ends, which closes the
// merged stream and lets runWatch return.
func TestRunWatch_StreamsEveryNamespace(t *testing.T) {
	kubeconfig := fakeAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/") // /api/v1/namespaces/<ns>/events
		if r.URL.Query().Get("watch") != "true" || len(parts) != 6 {
			http.NotFound(w, r)
			return
		}
		ns := parts[4]
		ev := k8sEvent("Pod", "pod-"+ns, ns, "BackOff", "Warning")
		ev.APIVersion, ev.Kind = "v1", "Event"
		obj, err := json.Marshal(ev)
		if err != nil {
			t.Error(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"type":"ADDED","object":%s}`+"\n", obj)
	})

	out, read := captureFile(t)
	f := eventFlags{since: "1h", output: "plain", namespaces: []string{"a", "b"}, kubeconfig: kubeconfig}
	if err := runWatch(f, out); err != nil {
		t.Fatalf("runWatch() error = %v", err)
	}

	got := read()
	for _, want := range []string{"Pod/pod-a [a]", "Pod/pod-b [b]"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
}

func TestRunWatch_ReportsARejectedWatch(t *testing.T) {
	kubeconfig := fakeAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	})

	err := runWatch(eventFlags{since: "1h", namespaces: []string{"a"}, kubeconfig: kubeconfig}, nil)
	if err == nil {
		t.Fatal("runWatch() error = nil, want the rejected watch reported")
	}
	if !strings.Contains(err.Error(), `namespace "a"`) {
		t.Errorf("error = %q, want it to name the namespace", err.Error())
	}
}
