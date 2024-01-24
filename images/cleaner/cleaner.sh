#!/bin/sh

echo "Running Connectivity Proxy Cleanup script"
kubectl get pods -n kyma-system

kubectl delete statefulset connectivity-proxy -n kyma-system

exit 0