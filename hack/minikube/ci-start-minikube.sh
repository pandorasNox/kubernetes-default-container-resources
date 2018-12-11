#!/bin/bash

export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export CHANGE_MINIKUBE_NONE_USER=true
mkdir -p "$HOME/.kube"
touch "$HOME/.kube/config"

export KUBECONFIG="$HOME/.kube/config"
sudo -E minikube start --vm-driver=none --extra-config=apiserver.admission-control="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,ValidatingAdmissionWebhook,MutatingAdmissionWebhook"
