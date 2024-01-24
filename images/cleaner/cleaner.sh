#!/bin/sh

echo "Running Connectivity Proxy Cleanup script"
kubectl get pods -n kyma-system
exit 0