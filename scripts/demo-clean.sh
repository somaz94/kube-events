#!/usr/bin/env bash
set -euo pipefail

NS="kube-events-demo"

echo "Cleaning up demo resources..."

kubectl delete namespace "$NS" --ignore-not-found --wait=false

echo "Done! Namespace '$NS' is being deleted."
