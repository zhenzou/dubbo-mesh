#!/bin/sh

root_path=$GOPATH/src/z/
cd $root_path
ignoreDir="(common/fs|app|tmp|build)"
go test $(go list ./... | grep -vE '(common/fs|app|tmp|build)')