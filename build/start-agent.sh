#!/bin/bash

ETCD_HOST=etcd
ETCD_PORT=2379
ETCD_URL=http://${ETCD_HOST}:${ETCD_PORT}

echo ETCD_URL = ${ETCD_URL}

export GOGC=400
#export GODEBUG=gctrace=1do

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  /root/dists/consumer -e=${ETCD_URL}
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  /root/dists/provider -m=2048 -ps=50 -n=provider-small -p=30000 -dp=20880 -e=${ETCD_URL}
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  /root/dists/provider -m=4096 -ps=100 -n=provider-medium -p=30000 -dp=20880 -e=${ETCD_URL}
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  /root/dists/provider -m=6144 -ps=150 -n=provider-large -p=30000 -dp=20880  -e=${ETCD_URL}
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
