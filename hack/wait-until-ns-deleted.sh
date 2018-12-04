#!/bin/bash

until [ $(kubectl get ns mutating-webhook 2>/dev/null | wc -l) -eq 0 ]; do
    echo -n "."
    sleep 1
done
echo "done"
