#!/usr/bin/env bash

docker run -d --name=provider-small --cpu-period=50000 --cpu-quota=30000  --memory=2g  --network=benchmarker dubbo-mesh provider-small
docker run -d --name=provider-medium --cpu-period=50000 --cpu-quota=60000 --memory=4g --network=benchmarker  dubbo-mesh provider-medium
docker run -d --name=provider-large  --cpu-period=50000 --cpu-quota=90000 --memory=6g --network=benchmarker  dubbo-mesh provider-large
docker run -d -p 8087:8087 --cpu-period 50000 --cpu-quota 180000 --memory=3g  --network=benchmarker  --name consumer dubbo-mesh consumer