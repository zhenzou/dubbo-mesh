#!/bin/bash

ETCD_HOST=etcd
ETCD_PORT=2379
ETCD_URL=http://${ETCD_HOST}:${ETCD_PORT}

#export GODEBUG=gctrace=1

echo ETCD_URL = ${ETCD_URL}

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  GOGC=50 GODEBUG=gctrace=1 /root/dists/consumer -e=${ETCD_URL}
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  /root/dists/provider -m=2048 -n=provider-small -p=30000 -dp=20880 -e=${ETCD_URL}
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  /root/dists/provider -m=4096 -n=provider-medium -p=30000 -dp=20880 -e=${ETCD_URL}
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  /root/dists/provider -m=6144 -n=provider-large -p=30000 -dp=20880  -e=${ETCD_URL}
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
