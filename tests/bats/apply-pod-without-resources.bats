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
    echo "debug | namespace: $test_ns" >&2
    kubectl create ns $test_ns > /dev/null

    cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: pod-without-resources
  namespace: $test_ns
spec:
  containers:
  - image: nginx:1.14.2-alpine
    name: nginx
    ports:
    - containerPort: 80
      protocol: TCP
EOF

    #wait for pod-without-resources is running
    until [ $(kubectl -n $test_ns get po 2> /dev/null | grep Running | wc -l) -eq 1 ]; do
        sleep 1
    done

    echo "debug | running pod: $(kubectl -n $test_ns get po | grep Running)" >&2

    resources_found=$(kubectl -n $test_ns get po -o yaml | grep -A 6 "resources:")

    echo "debug | resources_found: $resources_found" >&2

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

    echo "debug | expected_resources: $expected_resources" >&2

    # clean-up test namespace
    kubectl delete ns $test_ns > /dev/null || true

    result="0"
    if [ "$resources_found" == "$expected_resources" ]; then
      # echo "Strings match"
      result="1"
      echo "debug | resources_found == expected_resources" >&2
    else
      echo "debug | resources_found != expected_resources" >&2
      echo ""
    fi

    echo "debug | result: $result" >&2

    [ "$result" -eq 1 ]
}
