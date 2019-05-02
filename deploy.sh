#!/bin/bash
#
#
# Super Jank Deploy Script Courtesy De Nova

NAMESPACE="spark"
IMAGE="krisnova/spark-cluster-api-operator:latest"


make container push

kubectl delete namespace ${NAMESPACE}
kubectl create namespace ${NAMESPACE}
kubectl run spark-cluster-api-operator -n ${NAMESPACE} --image ${IMAGE} --env "KUBECONFIG_CONTENT=$(cat ~/.kube/config)"