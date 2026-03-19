#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY="${PROJECT_DIR}/kube-events"
NS="kube-events-demo"

# Colors
GREEN='\033[32m'
YELLOW='\033[33m'
CYAN='\033[36m'
BOLD='\033[1m'
RESET='\033[0m'

header() {
  echo ""
  echo -e "${BOLD}${CYAN}=== Phase $1: $2 ===${RESET}"
  echo ""
}

run() {
  echo -e "${YELLOW}\$ $*${RESET}"
  eval "$@"
  echo ""
}

# Build if needed
if [ ! -f "$BINARY" ]; then
  echo "Building kube-events..."
  (cd "$PROJECT_DIR" && make build)
fi

# ============================================================
header 1 "Deploy demo resources"
# ============================================================

run kubectl apply -f "${PROJECT_DIR}/examples/namespace.yaml"
run kubectl apply -f "${PROJECT_DIR}/examples/configmap.yaml"
run kubectl apply -f "${PROJECT_DIR}/examples/service.yaml"
run kubectl apply -f "${PROJECT_DIR}/examples/healthy-deployment.yaml"

echo "Waiting for healthy pods to be ready..."
kubectl wait --for=condition=available deployment/healthy-app -n "$NS" --timeout=60s || true
echo ""

# ============================================================
header 2 "Check events (should show Normal events)"
# ============================================================

run "$BINARY" -n "$NS" --since 5m

# ============================================================
header 3 "Deploy problematic resources"
# ============================================================

run kubectl apply -f "${PROJECT_DIR}/examples/crashloop-pod.yaml"
run kubectl apply -f "${PROJECT_DIR}/examples/bad-probe-deployment.yaml"
run kubectl apply -f "${PROJECT_DIR}/examples/bad-image-pod.yaml"

echo "Waiting 30s for events to accumulate..."
sleep 30

# ============================================================
header 4 "Check events (should show Warning events)"
# ============================================================

run "$BINARY" -n "$NS" --since 5m

# ============================================================
header 5 "Warning events only"
# ============================================================

run "$BINARY" -n "$NS" -t Warning --since 5m

# ============================================================
header 6 "Filter by kind (Pod only)"
# ============================================================

run "$BINARY" -n "$NS" -k Pod --since 5m

# ============================================================
header 7 "Filter by reason (BackOff)"
# ============================================================

run "$BINARY" -n "$NS" -r BackOff --since 5m

# ============================================================
header 8 "Filter by name"
# ============================================================

run "$BINARY" -n "$NS" -N crashloop-app --since 5m

# ============================================================
header 9 "Summary only"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -s

# ============================================================
header 10 "Group by namespace"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -g namespace

# ============================================================
header 11 "Group by kind"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -g kind

# ============================================================
header 12 "Group by reason"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -g reason

# ============================================================
header 13 "JSON output"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -o json

# ============================================================
header 14 "Markdown output"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -o markdown

# ============================================================
header 15 "Table output"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -o table

# ============================================================
header 16 "Plain output"
# ============================================================

run "$BINARY" -n "$NS" --since 5m -o plain

echo -e "${GREEN}${BOLD}Demo complete!${RESET}"
echo -e "Run ${CYAN}make demo-clean${RESET} to remove demo resources."
