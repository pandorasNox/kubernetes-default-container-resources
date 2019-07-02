#!/bin/bash

NAMESPACE=mutating-webhook

until [ $(kubectl get ns ${NAMESPACE} 2>/dev/null | wc -l) -eq 0 ]; do
    echo -n "."
    sleep 1
done
echo "done"
