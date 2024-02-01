#!/bin/sh

#Backups connectivity-proxy, and connectivity-proxy-info config map.
set -e

BACKUP_NS=cp-backup

echo "Running Connectivity Proxy Backup script"

echo "Checking config maps for backup"
#check if connectivity-proxy and connectivity-proxy-info config maps exists in namespace kyma-system and if not, exit with 0
if ! kubectl get cm -n kyma-system connectivity-proxy &> /dev/null; then
  echo "connectivity-proxy config map does not exist in kyma-system namespace, exiting"
  exit 0
fi

if ! kubectl get cm -n kyma-system connectivity-proxy-info &> /dev/null; then
  echo "connectivity-proxy-info config map does not exist in kyma-system namespace, exiting"
  exit 0
fi

# check if namespace from variable BACKUP_NS exists and if not create it in the same cluster
if ! kubectl get namespace $BACKUP_NS &> /dev/null; then
  echo "Namespace $BACKUP_NS does not exist, creating it"
  kubectl create namespace $BACKUP_NS
fi

#ensure that old backup is deleted before creating a new one
kubectl delete cm -n $BACKUP_NS connectivity-proxy --ignore-not-found
kubectl delete cm -n $BACKUP_NS connectivity-proxy-info --ignore-not-found

echo "Backing up Connectivity Proxy config maps"

# read config maps connectivity-proxy and connectivity-proxy-info config maps from kyma-system namespace and store them in a backup namespace
kubectl get cm -n kyma-system connectivity-proxy -o yaml | sed s/"namespace: kyma-system"/"namespace: $BACKUP_NS"/ | kubectl apply -f -
kubectl get cm -n kyma-system connectivity-proxy-info -o yaml | sed s/"namespace: kyma-system"/"namespace: $BACKUP_NS"/ | kubectl apply -f -


echo "Connectivity Proxy Backup script completed successfully"