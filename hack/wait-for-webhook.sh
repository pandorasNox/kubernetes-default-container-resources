#!/bin/bash

NAMESPACE=mutating-webhook
APP_SELECCTOR=app=default-container-resources

echo "wait until mutating-webhook is running"
until [ $(kubectl -n ${NAMESPACE} get po --selector=${APP_SELECCTOR} 2> /dev/null | grep Running | wc -l) -eq 1 ]; do
    echo -n "."
    sleep 1
done
echo "mutating-webhook is now running"

echo "wait until mutating-webhook is ready"
until [ $(kubectl -n ${NAMESPACE} get po --selector=${APP_SELECCTOR} 2> /dev/null | grep "1/1" | wc -l) -eq 1 ]; do
    echo -n "."
    sleep 1
done
echo "mutating-webhook is now ready"
