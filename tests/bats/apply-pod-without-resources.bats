#!/usr/bin/env bats

@test "admission controller is deployed" {
    result="$(kubectl -n mutating-webhook get po | wc -l)"
    [ "$result" -eq 2 ]
}

@test "admission controller is running" {
    result="$(kubectl -n mutating-webhook get po | grep Running | wc -l)"
    [ "$result" -eq 1 ]
}

@test "apply pod without resources sets default resources" {
    # ensure test namespace is not running
    kubectl delete ns bats-test > /dev/null || true
    until [ $(kubectl get ns bats-test 2> /dev/null | wc -l) -eq 0 ]; do
        sleep 1
    done

    # create test namespace
    kubectl create ns bats-test > /dev/null

    cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: pod-without-resources
  namespace: bats-test
spec:
  containers:
  - image: nginx:1.7.9
    name: nginx
    ports:
    - containerPort: 80
      protocol: TCP
EOF

    resources_found=$(kubectl -n bats-test get po -o yaml | grep -A 6 "resources:")

    expected_resources=$(cat <<-END
      resources:
        limits:
          cpu: 500m
          memory: 1G
        requests:
          cpu: 50m
          memory: 512M
END
)

    # clean-up test namespace
    kubectl delete ns bats-test > /dev/null || true

    result="0"
    if [ "$resources_found" == "$expected_resources" ]; then
      # echo "Strings match"
      result="1"
    fi

    [ "$result" -eq 1 ]
}
