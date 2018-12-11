#!/bin/bash

export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export CHANGE_MINIKUBE_NONE_USER=true
export KUBECONFIG=$HOME/.kube/config


# this for loop waits until kubectl can access the api server that Minikube has created
echo -n wait for minikube to start
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl version | grep "Server Version" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for local node to join
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl get no | grep " Ready " > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for dns to start
until [ $(kubectl -n kube-system -l k8s-app=kube-dns --no-headers=true get po | grep "Running" | wc -l) -eq 2 ]; do
    sleep 1
done
echo " done"

echo -n "wait for hole kube-system to be ready"
until [ $(kubectl -n kube-system get po --no-headers=true | grep -v "Running" | wc -l) -eq 0 ]; do
    sleep 1
done
echo " done"

# kubectl commands are now able to interact with Minikube cluster

# workaround https://github.com/kubernetes/minikube/issues/1947
# echo -n getting name of kubedns pod
# for _ in {1..150}; do # timeout for 5 minutes
#   if KUBEDNS_POD=$(kubectl get --no-headers=true pods -n kube-system -l k8s-app=kube-dns -o custom-columns=:metadata.name); then
#       echo " done"
#       break
#   fi
#   echo -n .
#   sleep 2
# done

# echo -n fixing kubedns upstream server
# for _ in {1..150}; do # timeout for 5 minutes
#   if kubectl exec -n kube-system "$KUBEDNS_POD" -c kubedns -- sh -c "echo nameserver 8.8.8.8 > /etc/resolv.conf"; then
#       echo " done"
#       break
#   fi
#   echo -n .
#   sleep 2
# done
