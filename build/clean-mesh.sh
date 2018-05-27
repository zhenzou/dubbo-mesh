#!/usr/bin/env bash
echo "clean provider-small"
docker stop provider-small
docker rm provider-small
echo "clean provider-medium"
docker stop provider-medium
docker rm provider-medium
echo "clean provider-medium"
docker stop provider-large
docker rm provider-large
echo "clean consumer"
docker stop consumer
docker rm consumer
