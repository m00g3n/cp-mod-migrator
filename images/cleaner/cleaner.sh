#!/bin/sh

echo "Running Connectivity Proxy Cleanup script"

kubectl delete statefulset -n kyma-system connectivity-proxy
kubectl delete deployment -n kyma-system connectivity-proxy-restart-watcher
kubectl delete deployment -n kyma-system connectivity-proxy-sm-operator

kubectl delete service -n kyma-system connectivity-proxy
kubectl delete service -n kyma-system connectivity-proxy-smv
kubectl delete service -n kyma-system connectivity-proxy-tunnel
kubectl delete service -n kyma-system connectivity-proxy-tunnel-0
kubectl delete service -n kyma-system connectivity-proxy-tunnel-healthcheck

kubectl delete serviceaccount -n kyma-system connectivity-proxy-restart-watcher
kubectl delete clusterrole -n kyma-system connectivity-proxy-restart-watcher
kubectl delete clusterrolebinding -n kyma-system connectivity-proxy-restart-watcher

kubectl delete serviceaccount -n kyma-system connectivity-proxy-sm-operator
kubectl delete clusterrole -n kyma-system connectivity-proxy-service-mappings
kubectl delete clusterrolebinding -n kyma-system connectivity-proxy-service-mappings

kubectl delete configmap -n kyma-system connectivity-proxy
kubectl delete configmap -n kyma-system connectivity-proxy-info


exit 0