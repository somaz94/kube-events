#!/usr/bin/env bash
#
# End-to-end check for watch mode against a live cluster.
#
# It guards the behavior unit tests cannot reach: --namespace is repeatable, so
# `--watch` must open a watch on EVERY namespace given, not just the first. The
# check starts a watcher across two namespaces, makes an event happen in each,
# and fails unless both surface on the stream.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY="${PROJECT_DIR}/kube-events"

NS_A="kube-events-watch-a"
NS_B="kube-events-watch-b"
POD_A="watch-probe-a"
POD_B="watch-probe-b"

# How long to wait for both events to reach the stream.
TIMEOUT_SECONDS="${WATCH_TIMEOUT_SECONDS:-90}"

OUT_FILE="$(mktemp)"
WATCH_PID=""

# shellcheck disable=SC2329  # invoked through the EXIT trap below
cleanup() {
  if [ -n "$WATCH_PID" ] && kill -0 "$WATCH_PID" 2>/dev/null; then
    kill "$WATCH_PID" 2>/dev/null || true
    wait "$WATCH_PID" 2>/dev/null || true
  fi
  kubectl delete namespace "$NS_A" "$NS_B" --ignore-not-found --wait=false >/dev/null 2>&1 || true
  rm -f "$OUT_FILE"
}
trap cleanup EXIT

fail() {
  echo "FAIL: $*" >&2
  echo "--- watch output ---" >&2
  cat "$OUT_FILE" >&2
  exit 1
}

# This check creates and deletes namespaces, so it must not be pointed at a real
# cluster by accident. Only a kind context is allowed unless the caller opts out
# deliberately.
CONTEXT="$(kubectl config current-context 2>/dev/null || echo '')"
case "$CONTEXT" in
  kind-*) ;;
  *)
    if [ "${ALLOW_ANY_CONTEXT:-0}" != "1" ]; then
      echo "refusing to run against context '${CONTEXT:-<none>}': this check creates and deletes namespaces." >&2
      echo "Use a kind cluster, or set ALLOW_ANY_CONTEXT=1 to override." >&2
      exit 1
    fi
    echo "WARNING: running against non-kind context '$CONTEXT' (ALLOW_ANY_CONTEXT=1)"
    ;;
esac

if [ ! -x "$BINARY" ]; then
  echo "Building kube-events..."
  (cd "$PROJECT_DIR" && make build)
fi

echo "==> Creating namespaces"
kubectl create namespace "$NS_A" >/dev/null
kubectl create namespace "$NS_B" >/dev/null

# The watcher must be running before the events happen, otherwise there is
# nothing live to observe.
echo "==> Starting watch across both namespaces"
"$BINARY" -n "$NS_A" -n "$NS_B" --watch --since 5m > "$OUT_FILE" 2>&1 &
WATCH_PID=$!

# Give the watch connections time to establish.
sleep 5

if ! kill -0 "$WATCH_PID" 2>/dev/null; then
  fail "the watcher exited before any event was produced"
fi

echo "==> Producing an event in each namespace"
kubectl run "$POD_A" --image=busybox:1.36 --restart=Never -n "$NS_A" -- sleep 30 >/dev/null
kubectl run "$POD_B" --image=busybox:1.36 --restart=Never -n "$NS_B" -- sleep 30 >/dev/null

echo "==> Waiting for both namespaces to appear on the stream (timeout ${TIMEOUT_SECONDS}s)"
deadline=$((SECONDS + TIMEOUT_SECONDS))
while [ "$SECONDS" -lt "$deadline" ]; do
  if grep -q "$POD_A" "$OUT_FILE" && grep -q "$POD_B" "$OUT_FILE"; then
    echo "PASS: events from both namespaces reached the watch stream"
    exit 0
  fi
  if ! kill -0 "$WATCH_PID" 2>/dev/null; then
    fail "the watcher exited early"
  fi
  sleep 2
done

# Report precisely which side is missing: seeing only the first namespace is the
# exact symptom of the single-watch regression this check exists for.
saw_a=no; saw_b=no
grep -q "$POD_A" "$OUT_FILE" && saw_a=yes
grep -q "$POD_B" "$OUT_FILE" && saw_b=yes
fail "timed out waiting for events (saw ${NS_A}=${saw_a}, ${NS_B}=${saw_b})"
