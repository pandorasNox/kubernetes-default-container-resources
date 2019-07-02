#!/bin/bash

NAMESPACE=mutating-webhook
APP_SELECCTOR=app=default-container-resources
EXPECTED_RUNNING=1
EXPECTED_READY=1/1

echo "wait until mutating-webhook is running"
until [ $(kubectl -n ${NAMESPACE} get po --selector=${APP_SELECCTOR} 2> /dev/null | grep Running | wc -l) -eq ${EXPECTED_RUNNING} ]; do
    echo -n "."
    sleep 1
done
echo "mutating-webhook is now running"

echo "wait until mutating-webhook is ready"
until [ $(kubectl -n ${NAMESPACE} get po --selector=${APP_SELECCTOR} 2> /dev/null | grep "${EXPECTED_READY}" | wc -l) -eq ${EXPECTED_RUNNING} ]; do
    echo -n "."
    sleep 1
done
echo "mutating-webhook is now ready"
