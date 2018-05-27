#!/usr/bin/env bash

echo "restart provider-small"
docker restart provider-small
echo "restart provider-medium"
docker restart provider-medium
echo "restart provider-large"
docker restart provider-large
echo "restart consumer"
docker restart consumer
