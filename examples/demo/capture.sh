#!/usr/bin/env bash
# Runs the README demo for real on a throwaway kind cluster and saves the
# kube-events output under examples/demo/captured/, which replay.sh prints back.
# Usage: make build && examples/demo/capture.sh   (needs kind and kubectl; pulls nginx and busybox)
# bad-image-pod.yaml is left out: its pull error is one 374-character line.
set -euo pipefail

DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
ROOT_DIR="$(cd "${DEMO_DIR}/../.." && pwd)"
OUT="${DEMO_DIR}/captured"
CLUSTER="kube-events-demo"
NS="kube-events-demo"
SETTLE="${SETTLE:-60}"

for tool in kind kubectl; do
    command -v "${tool}" >/dev/null 2>&1 || { echo "${tool} not found" >&2; exit 1; }
done
[ -x "${ROOT_DIR}/kube-events" ] || { echo "kube-events not found; run make build" >&2; exit 1; }
if kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
    echo "kind cluster ${CLUSTER} already exists; delete it first" >&2
    exit 1
fi

WORK="$(mktemp -d)"
# kind would otherwise switch the current context of the user's own kubeconfig.
export KUBECONFIG="${WORK}/kubeconfig"
cleanup() {
    kind delete cluster --name "${CLUSTER}" --kubeconfig "${KUBECONFIG}" >/dev/null 2>&1 || true
    rm -rf "${WORK:?}"
}
trap cleanup EXIT

kind create cluster --name "${CLUSTER}" --kubeconfig "${KUBECONFIG}" --wait 60s

cd "${ROOT_DIR}"
kubectl apply -f examples/namespace.yaml
kubectl apply -f examples/configmap.yaml -f examples/service.yaml -f examples/healthy-deployment.yaml \
    -f examples/crashloop-pod.yaml -f examples/bad-probe-deployment.yaml
kubectl wait --for=condition=available deployment/healthy-app -n "${NS}" --timeout=120s || true
echo "waiting ${SETTLE}s for events to accumulate"
sleep "${SETTLE}"

rm -rf "${OUT}"
mkdir -p "${OUT}"
./kube-events -n "${NS}" -k Pod --since 5m > "${OUT}/01-kube-events" 2>&1 || true

echo "captured into ${OUT}"
