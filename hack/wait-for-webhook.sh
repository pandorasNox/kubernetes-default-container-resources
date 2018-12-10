#!/bin/bash

until [ $(kubectl -n mutating-webhook get po 2> /dev/null | grep Running | wc -l) -eq 1 ]; do
    echo -n "."
    sleep 1
done
echo "mutating-webhook is now running"
