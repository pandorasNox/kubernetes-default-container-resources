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

    # create test namespace
    test_ns=draco-test-$(uuidgen)
    kubectl create ns $test_ns > /dev/null

    cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: pod-without-resources
  namespace: $test_ns
spec:
  containers:
  - image: nginx:1.7.9
    name: nginx
    ports:
    - containerPort: 80
      protocol: TCP
EOF

    resources_found=$(kubectl -n $test_ns get po -o yaml | grep -A 6 "resources:")

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
    kubectl delete ns $test_ns > /dev/null || true

    result="0"
    if [ "$resources_found" == "$expected_resources" ]; then
      # echo "Strings match"
      result="1"
    fi

    [ "$result" -eq 1 ]
}
