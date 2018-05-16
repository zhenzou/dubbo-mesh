#!/bin/bash

ETCD_HOST=$(ip addr show docker0 | grep 'inet\b' | awk '{print $2}' | cut -d '/' -f 1)
ETCD_PORT=2379
ETCD_URL=http://${ETCD_HOST}:${ETCD_PORT}

echo ETCD_URL = ${ETCD_URL}

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  /root/dists/consumer -e=${ETCD_URL}
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  /root/dists/provider -m=2048 -n=provider-small -p=30000 -dp=20889 -e=${ETCD_URL}
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  /root/dists/provider -m=4096 -n=provider-medium -p=30001 -dp=20890 -e=${ETCD_URL}
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  /root/dists/provider -m=6144 -n=provider-large -p=30002 -dp=20891  -e=${ETCD_URL}
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
