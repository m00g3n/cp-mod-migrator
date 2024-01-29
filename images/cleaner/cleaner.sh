#!/bin/sh

set -e

echo "Running Connectivity Proxy Cleanup script"

# check if statefulset connectivity-proxy does not exist and if so, exit  with 0
if ! kubectl get statefulset -n kyma-system connectivity-proxy &> /dev/null; then
  echo "Connectivity Proxy does not exist, exiting"
  exit 0
fi

echo "Connectivity Proxy is installed"

echo "Annotate all existing Connectivity Proxy Service Mappings"

mappings=$(kubectl get servicemappings.connectivityproxy.sap.com -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')

for mapping in $mappings; do

  echo "Applying annotations to service mapping $mapping"

  kubectl annotate servicemappings.connectivityproxy.sap.com "$mapping" \
    io.javaoperatorsdk/primary-name=connectivity-proxy \
    io.javaoperatorsdk/primary-namespace=kyma-system \

done

echo "Removing all Connectivity Proxy workloads"

kubectl delete statefulset -n kyma-system connectivity-proxy --ignore-not-found
kubectl delete deployment -n kyma-system connectivity-proxy-restart-watcher --ignore-not-found
kubectl delete deployment -n kyma-system connectivity-proxy-sm-operator --ignore-not-found

echo "Removing all Connectivity Proxy services"

kubectl delete service -n kyma-system connectivity-proxy --ignore-not-found
kubectl delete service -n kyma-system connectivity-proxy-smv --ignore-not-found
kubectl delete service -n kyma-system connectivity-proxy-tunnel --ignore-not-found
kubectl delete service -n kyma-system connectivity-proxy-tunnel-0 --ignore-not-found
kubectl delete service -n kyma-system connectivity-proxy-tunnel-healthcheck --ignore-not-found

echo "Removing all Connectivity Proxy RBAC resources"

kubectl delete clusterrolebinding connectivity-proxy-restart-watcher --ignore-not-found
kubectl delete clusterrole connectivity-proxy-restart-watcher --ignore-not-found
kubectl delete serviceaccount -n kyma-system connectivity-proxy-restart-watcher --ignore-not-found

kubectl delete clusterrolebinding connectivity-proxy-service-mappings --ignore-not-found
kubectl delete clusterrole connectivity-proxy-service-mappings --ignore-not-found
kubectl delete serviceaccount -n kyma-system connectivity-proxy-sm-operator --ignore-not-found

kubectl delete mutatingwebhookconfiguration connectivity-proxy-mutating-webhook-configuration --ignore-not-found

echo "Removing all Connectivity Proxy ConfigMaps"

kubectl delete configmap -n kyma-system connectivity-proxy --ignore-not-found
kubectl delete configmap -n kyma-system connectivity-proxy-info --ignore-not-found

echo "Removing all Connectivity Proxy Istio resources"

kubectl delete envoyfilter -n istio-system connectivity-proxy-custom-protocol --ignore-not-found
kubectl delete gateway -n kyma-system kyma-gateway-cc --ignore-not-found
kubectl delete virtualservice -n kyma-system cc-proxy --ignore-not-found
kubectl delete virtualservice -n kyma-system cc-proxy-healthcheck --ignore-not-found
kubectl delete destinationrule -n kyma-system connectivity-proxy-tunnel-0 --ignore-not-found
kubectl delete destinationrule -n kyma-system connectivity-proxy --ignore-not-found
kubectl delete peerauthentication -n kyma-system enable-permissive-mode-for-cp --ignore-not-found
kubectl delete certificate -n istio-system cc-certs --ignore-not-found
kubectl delete secret -n istio-system cc-certs-cacert --ignore-not-found

