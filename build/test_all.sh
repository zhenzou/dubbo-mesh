#!/bin/sh

PROJ_ROOT=`dirname $0`
cd "${PROJ_ROOT}/../"
ignoreDir="(common/fs|app|tmp|build)"
go test $(go list ./... | grep -vE '(common/fs|app|tmp|build)')