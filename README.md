# kubernetes-default-container-resources
a Kubernetes admission controller mutating webhook which adds default compute resources to container

### generate cert.pem and key.pem for ListenAndServeTLS
`docker run --rm -v $(pwd)/certs:/certs -e SSL_SUBJECT=test.example.com -e SSL_KEY="ssl-key.pem" -e SSL_CSR="ssl-key.csr" -e SSL_CERT="ssl-cert.pem" -e K8S_NAME="pls-replace-me-kubernetes-name" paulczar/omgwtfssl`

after that, to change permissions ... run:
`sudo chown -R username:groupname certs`


###
without resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

with empty resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

with memory 'limit' resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"limits":{"memory": "512Mi"}},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

with cpu 'limit' resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"limits":{"cpu": "0.5"}},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

with memory 'requests' resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"requests":{"memory": "512Mi"}},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

with cpu 'requests' resources
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"requests":{"cpu": "0.5"}},"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}' \
  http://localhost:8083/mutate
```

_________________________________________

pod without resources
```
{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"83b928e0-89b3-11e8-b437-08002738f604","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"minikube-user","groups":["system:masters","system:authenticated"]},"object":{"metadata":{"name":"pod-without-resources","namespace":"foo","creationTimestamp":null,"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"pod-without-resources\",\"namespace\":\"foo\"},\"spec\":{\"containers\":[{\"image\":\"nginx:1.7.9\",\"name\":\"nginx\",\"ports\":[{\"containerPort\":80,\"protocol\":\"TCP\"}]}]}}\n"}},"spec":{"volumes":[{"name":"default-token-crw24","secret":{"secretName":"default-token-crw24"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{},"volumeMounts":[{"name":"default-token-crw24","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{}},"oldObject":null}}
```


###
- reSorcerer (Resource)

### usefull tools
- https://json-patch-builder-online.github.io/
